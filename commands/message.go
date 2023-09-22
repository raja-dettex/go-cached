package commands

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Command string

const (
	CMDSet    Command = "SET"
	CMDGet    Command = "GET"
	CMDHas    Command = "HAS"
	CMDDelete Command = "DELETE"
)

type Msg struct {
	Cmd   Command
	Key   []byte
	Value []byte
	TTL   time.Duration
}

func ParseMessage(raw []byte) (*Msg, error) {
	command := strings.Split(string(raw), " ")
	if len(command) < 2 {
		return nil, errors.New("invalid command")
	}
	msg := &Msg{
		Cmd: Command(command[0]),
		Key: []byte(command[1]),
	}
	if msg.Cmd == CMDSet {
		if len(command) < 4 {
			return nil, errors.New("invalid Set Command")
		}
		msg.Value = []byte(command[2])
		ttl, err := strconv.Atoi(command[3])
		if err != nil {
			return nil, errors.New("invalid set TTL")
		}
		msg.TTL = time.Duration(ttl)
		return msg, nil
	}
	return msg, nil
}
