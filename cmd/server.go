package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/sir-hassan/bambus/auth"
	backendFactory "github.com/sir-hassan/bambus/backend/factory"
	"github.com/sir-hassan/bambus/config"
	"github.com/sir-hassan/bambus/core"
	"github.com/sir-hassan/bambus/frontend"
	socketFactory "github.com/sir-hassan/bambus/frontend/factory"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const wssPort = 8080

func main() {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	level.Info(logger).Log("msg", "starting bambus...")

	configBytes, err := os.ReadFile("./config.json")
	if err != nil {
		level.Error(logger).Log("msg", "opening config file", "error", err)
		return
	}
	var appConfig config.AppConfig
	if err := json.Unmarshal(configBytes, &appConfig); err != nil {
		level.Error(logger).Log("msg", "parsing config file", "error", err)
		return
	}
	tubeCreator, err := backendFactory.CreateTubeCreator(logger, appConfig.Backend)
	if err != nil {
		level.Error(logger).Log("msg", "pubSub backend", "error", err)
		return
	}
	socketType, err := appConfig.Frontend.GetStringField(
		"type",
		config.NoEmptyString,
	)
	if err != nil {
		level.Error(logger).Log("msg", "frontend config", "error", err)
		return
	}
	socketCreator, err := socketFactory.CreateSocketCreator(socketType)
	if err != nil {
		level.Error(logger).Log("msg", "frontend", "error", err)
		return
	}
	hub := core.NewHub(logger, tubeCreator)
	authenticator, err := auth.AuthenticatorCreator(appConfig.Authenticator)
	if err != nil {
		level.Error(logger).Log("msg", "authenticator config", "error", err)
		return
	}

	http.HandleFunc("/wss/", createConnHandler(logger, hub, authenticator, socketCreator))
	level.Info(logger).Log("msg", "listening started", "port", wssPort)
	ln, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(wssPort))
	if err != nil {
		level.Error(logger).Log("msg", "listening on port failed", "port", wssPort, "error", err)
		return
	}
	server := &http.Server{}

	go func() {
		err := server.Serve(ln)
		if err != http.ErrServerClosed {
			level.Error(logger).Log("msg", "closing server failed", "error", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	level.Info(logger).Log("msg", "signal received", "sig", sig)
	level.Info(logger).Log("msg", "terminating server...")

	if err := server.Shutdown(context.Background()); err != nil {
		level.Error(logger).Log("msg", "terminating server failed", "error", err)
	} else {
		level.Info(logger).Log("msg", "terminating server succeeded")
	}
}

func createConnHandler(logger log.Logger, hub *core.Hub, auth auth.Authenticator, socketCreator frontend.SocketCreator) func(w http.ResponseWriter, r *http.Request) {
	index := 0
	return func(w http.ResponseWriter, r *http.Request) {

		channels, err := auth.Auth(r)
		if err != nil {
			level.Error(logger).Log("msg", "authenticator", "error", err)
			writeReply(logger, w, 500, "internal server error")
			return
		}
		if len(channels) == 0 {
			writeReply(logger, w, 401, "unauthorized")
			return
		}

		soc, err := socketCreator(w, r)
		if err != nil {
			level.Error(logger).Log("msg", "socket create", "error", err)
			writeReply(logger, w, 500, "internal server error")
			return
		}
		level.Info(logger).Log("msg", "connection opened")
		index = index + 1
		unPlug := hub.Plug(soc, channels)

		go func() {
			defer func() {
				unPlug()
				//level.Error(logger).Log("msg", "remove socket", "error", err)
				level.Info(logger).Log("msg", "connection closed")
			}()
			for {
				message, err := soc.Read()
				if err != nil {
					level.Error(logger).Log("msg", "read socket", "error", err)
					return
				}
				level.Debug(logger).Log("msg", "read socket", "data", message)
				if err = soc.Write(message); err != nil {
					level.Error(logger).Log("msg", "write socket", "error", err)
					return
				}
			}
		}()
	}
}

func writeReply(logger log.Logger, w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		level.Error(logger).Log("msg", "writing connection", "error", err)
	}
}
