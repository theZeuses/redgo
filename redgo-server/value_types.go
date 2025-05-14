package main

type ValueType string

const (
	STRING      = '+'
	ERROR       = '-'
	INTEGER     = ':'
	BULK_STRING = '$'
	ARRAY       = '*'
	NULL        = '_'
)

const (
	R_STRING      ValueType = "string"
	R_ERROR       ValueType = "error"
	R_INTEGER     ValueType = "integer"
	R_BULK_STRING ValueType = "bulk_string"
	R_ARRAY       ValueType = "array"
	R_NULL        ValueType = "null"
	R_EMPTY       ValueType = "empty"
)

type Value interface {
	Type() ValueType
	Marshal() []byte
}

type StringValue struct {
	Val string
}

func (s StringValue) Type() ValueType { return R_STRING }

type ErrorValue struct {
	Val string
}

func (e ErrorValue) Type() ValueType { return R_ERROR }

type IntegerValue struct {
	Val int
}

func (i IntegerValue) Type() ValueType { return R_INTEGER }

type BulkStringValue struct {
	Val string
}

func (b BulkStringValue) Type() ValueType { return R_BULK_STRING }

type ArrayValue struct {
	Val []Value
}

func (a ArrayValue) Type() ValueType { return R_ARRAY }

type NullValue struct{}

func (n NullValue) Type() ValueType { return R_NULL }

type EmptyValue struct{}

func (e EmptyValue) Type() ValueType { return R_EMPTY }
