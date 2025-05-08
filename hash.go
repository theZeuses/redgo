package main

func hset(args []Value, _ *Client) Value {
	if len(args) < 3 || len(args)%2 != 1 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'hset' command"}
	}

	key := args[0].(BulkStringValue).Val

	HSETsMutex.Lock()
	defer HSETsMutex.Unlock()

	if _, found := HSETs[key]; !found {
		HSETs[key] = make(map[string]string)
	}

	for i := 1; i < len(args)-1; i += 2 {
		field := args[i].(BulkStringValue).Val
		value := args[i+1].(BulkStringValue).Val

		HSETs[key][field] = value
	}

	return StringValue{Val: "OK"}
}

func hget(args []Value, _ *Client) Value {
	if len(args) != 2 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'hget' command"}
	}

	key := args[0].(BulkStringValue).Val
	field := args[1].(BulkStringValue).Val

	HSETsMutex.RLock()
	defer HSETsMutex.RUnlock()

	if value, found := HSETs[key][field]; found {
		return BulkStringValue{Val: value}
	}

	return NullValue{}
}

func hgetall(args []Value, _ *Client) Value {
	if len(args) != 1 {
		return ErrorValue{Val: "ERR wrong number of arguments for 'hgetall' command"}
	}

	key := args[0].(BulkStringValue).Val

	HSETsMutex.RLock()
	hash, found := HSETs[key]
	HSETsMutex.RUnlock()

	if !found {
		return NullValue{}
	}

	array := make([]Value, 0, len(hash)*2)
	for field, value := range hash {
		array = append(array, BulkStringValue{Val: field})
		array = append(array, BulkStringValue{Val: value})
	}
	return ArrayValue{Val: array}
}
