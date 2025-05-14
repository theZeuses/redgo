package main

import (
	"net"
	"sync"
)

type Client struct {
	ID            string
	Conn          net.Conn
	Subscriptions map[string]*PubSubChannel
	Reader        *Reader
	Writer        *Writer
}

type PubSubChannelClient struct {
	Client *Client
	Next   *PubSubChannelClient
}

type PubSubChannelClientList struct {
	Head *PubSubChannelClient
}

type PubSubChannel struct {
	Name    string
	Clients PubSubChannelClientList
}

var (
	Clients             = make(map[string]*Client)
	PubSubChannels      = make(map[string]*PubSubChannel)
	ClientsMutex        sync.Mutex
	PubSubChannelsMutex sync.Mutex
)

func (list *PubSubChannelClientList) FindClientByID(id string) *Client {
	for client := list.Head; client != nil; client = client.Next {
		if client.Client.ID == id {
			return client.Client
		}
	}
	return nil
}

func (list *PubSubChannelClientList) AddClient(client *Client) {
	newClient := &PubSubChannelClient{Client: client}
	if list.Head == nil {
		list.Head = newClient
	} else {
		current := list.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newClient
	}
}

/**
 * RemoveClientByID removes a client from the list by its ID.
 * It returns true if the client was found and removed, and false if the client was not found.
 * The second return value indicates if the list is empty after the removal.
 */
func (list *PubSubChannelClientList) RemoveClientByID(id string) (bool, bool) {
	if list.Head == nil {
		return false, true
	}

	if list.Head.Client.ID == id {
		list.Head = list.Head.Next
		return true, true
	}

	current := list.Head
	for current.Next != nil {
		if current.Next.Client.ID == id {
			current.Next = current.Next.Next
			return true, false
		}
		current = current.Next
	}
	return false, false
}

func (list *PubSubChannelClientList) Len() int {
	count := 0
	for client := list.Head; client != nil; client = client.Next {
		count++
	}
	return count
}

func subscribe(args []Value, client *Client) Value {
	if len(args) < 1 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'subscribe' command"}
	}

	PubSubChannelsMutex.Lock()
	defer PubSubChannelsMutex.Unlock()

	for _, val := range args {
		channel := val.(BulkStringValue).Val
		if _, found := PubSubChannels[channel]; !found {
			PubSubChannels[channel] = &PubSubChannel{
				Name:    channel,
				Clients: PubSubChannelClientList{Head: nil},
			}
		}

		PubSubChannels[channel].Clients.AddClient(client)
		client.Subscriptions[channel] = PubSubChannels[channel]

		commandResponse := make([]Value, 3)
		commandResponse[0] = BulkStringValue{Val: "subscribe"}
		commandResponse[1] = BulkStringValue{Val: channel}
		commandResponse[2] = IntegerValue{Val: len(client.Subscriptions)}

		client.Writer.WriteAsRespString(ArrayValue{Val: commandResponse})
	}

	return EmptyValue{}
}

func unsubscribe(args []Value, client *Client) Value {
	if len(args) < 1 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'unsubscribe' command"}
	}

	PubSubChannelsMutex.Lock()
	defer PubSubChannelsMutex.Unlock()

	for _, val := range args {
		channel := val.(BulkStringValue).Val
		if pubsubChannel, found := PubSubChannels[channel]; found {
			_, isEmpty := pubsubChannel.Clients.RemoveClientByID(client.ID)
			delete(client.Subscriptions, channel)
			if isEmpty {
				delete(PubSubChannels, channel)
			}
		}
	}

	return StringValue{Val: "OK"}
}

func publish(args []Value, _ *Client) Value {
	if len(args) != 2 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'publish' command"}
	}

	PubSubChannelsMutex.Lock()
	defer PubSubChannelsMutex.Unlock()

	channel := args[0].(BulkStringValue).Val
	message := args[1].(BulkStringValue).Val

	pubsubChannel, found := PubSubChannels[channel]

	if !found {
		return IntegerValue{Val: 0}
	}

	head := pubsubChannel.Clients.Head
	len := 0

	for head != nil {
		client := head.Client
		subscriberResponse := make([]Value, 3)
		subscriberResponse[0] = BulkStringValue{Val: "message"}
		subscriberResponse[1] = BulkStringValue{Val: channel}
		subscriberResponse[2] = BulkStringValue{Val: message}

		client.Writer.WriteAsRespString(ArrayValue{Val: subscriberResponse})
		head = head.Next
		len++
	}

	return IntegerValue{Val: len}
}

func unsubscribeAll(client *Client) {
	PubSubChannelsMutex.Lock()
	defer PubSubChannelsMutex.Unlock()

	for channel := range client.Subscriptions {
		if pubsubChannel, found := PubSubChannels[channel]; found {
			_, isEmpty := pubsubChannel.Clients.RemoveClientByID(client.ID)
			if isEmpty {
				delete(PubSubChannels, channel)
			}
		}
	}

	client.Subscriptions = make(map[string]*PubSubChannel)
}
