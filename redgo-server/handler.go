package main

import (
	"strings"
	"sync"
)

var Handlers = map[string]func([]Value, *Client) Value{
	"PING":        ping,
	"SET":         set,
	"DEL":         del,
	"GET":         get,
	"HSET":        hset,
	"HGET":        hget,
	"HGETALL":     hgetall,
	"SUBSCRIBE":   subscribe,
	"UNSUBSCRIBE": unsubscribe,
	"PUBLISH":     publish,
}

var SETs = map[string]string{}
var HSETs = map[string]map[string]string{}

var SETsMutex = sync.RWMutex{}
var HSETsMutex = sync.RWMutex{}

func ping(args []Value, _ *Client) Value {
	if len(args) > 1 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'ping' command"}
	}

	msg := "PONG"

	if len(args) == 1 {
		msg = args[0].(BulkStringValue).Val
	}

	return BulkStringValue{Val: msg}
}

func ProcessCommand(command string, args []Value, client *Client) Value {
	if handler, found := Handlers[command]; found {
		return handler(args, client)
	}

	return ErrorValue{Val: "ERR unknown command '" + command + "'"}
}

func Handle(client Client, aof *Aof) error {
	for {
		value, err := client.Reader.ParseFromRespString()
		if err != nil {
			client.Writer.WriteAsRespString(ErrorValue{Val: "ERR parsing error"})
			return err
		}
		println("Received command from client:", client.ID)

		arrayVal, ok := value.(ArrayValue)

		if !ok || len(arrayVal.Val) == 0 {
			client.Writer.WriteAsRespString(ErrorValue{Val: "ERR wrong number of arguments"})
		} else {
			command := strings.ToUpper(arrayVal.Val[0].(BulkStringValue).Val)
			args := arrayVal.Val[1:]

			response := ProcessCommand(command, args, &client)

			if command == "SET" || command == "HSET" || command == "DEL" {
				aof.Write(value)
			}

			client.Writer.WriteAsRespString(response)
		}
	}
}
