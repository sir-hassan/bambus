package core

import (
	"fmt"
	"github.com/sir-hassan/bambus/frontend"
)

type contact struct {
	socketQueue frontend.SocketQueue
	channels    []string
	passUnPlug  chan func()
}

type contactsTable struct {
	channelsMap map[string][]frontend.SocketQueue
}

func newContactsTable() *contactsTable {
	return &contactsTable{
		channelsMap: make(map[string][]frontend.SocketQueue),
	}
}

func (s *contactsTable) Add(soc frontend.SocketQueue, channels []string) {
	for _, c := range channels {
		_, ok := s.channelsMap[c]
		if !ok {
			s.channelsMap[c] = make([]frontend.SocketQueue, 0, 0)
		}
		s.channelsMap[c] = append(s.channelsMap[c], soc)
	}
	printAllChannel(s)
}

func (s *contactsTable) Remove(soc frontend.SocketQueue) {
	for c, list := range s.channelsMap {
		list = removeSocketFromList(list, soc)
		s.channelsMap[c] = list
	}
	printAllChannel(s)
}

func (s *contactsTable) GetSocketsInChannel(channel string) []frontend.SocketQueue {
	if list, ok := s.channelsMap[channel]; ok {
		cp := make([]frontend.SocketQueue, len(list))
		copy(cp, list)
		return cp
	}
	return nil
}

func (s *contactsTable) GarbageCollect() []string {
	garbage := make([]string, 0)
	for channel, v := range s.channelsMap {
		if len(v) == 0 {
			garbage = append(garbage, channel)
			delete(s.channelsMap, channel)
		}
	}
	return garbage
}

func removeSocketFromList(list []frontend.SocketQueue, soc frontend.SocketQueue) []frontend.SocketQueue {
	for i, s := range list {
		if s == soc {
			list[i] = list[len(list)-1]
			list = list[:len(list)-1]
			return list
		}
	}
	return list
}

// todo: remove debug code.
func printAllChannel(s *contactsTable) {
	fmt.Println("----------------------------------------------------")
	for channel, _ := range s.channelsMap {
		socs := s.GetSocketsInChannel(channel)
		fmt.Printf("channel: %s: socs: %v\n", channel, socs)
	}
	fmt.Println("----------------------------------------------------")
}
