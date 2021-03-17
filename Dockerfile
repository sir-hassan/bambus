FROM golang:1.16-buster as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build

FROM scratch
COPY --chown=0:0 --from=builder /app/bin/bambus /app/bambus

# Todo: change to none root user
WORKDIR /app
ENTRYPOINT ["/app/bambus"]