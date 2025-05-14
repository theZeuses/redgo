package main

func set(args []Value, _ *Client) Value {
	if len(args) != 2 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].(BulkStringValue).Val
	value := args[1].(BulkStringValue).Val

	SETsMutex.Lock()
	defer SETsMutex.Unlock()

	SETs[key] = value

	return StringValue{Val: "OK"}
}

func get(args []Value, _ *Client) Value {
	if len(args) != 1 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].(BulkStringValue).Val

	SETsMutex.RLock()
	defer SETsMutex.RUnlock()

	if value, found := SETs[key]; found {
		return BulkStringValue{Val: value}
	}

	return NullValue{}
}
