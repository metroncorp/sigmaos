package mr

import (
	"bufio"
)

// Map and reduce functions produce and consume KeyValue pairs
type KeyValue struct {
	Key   string
	Value string
}

type EmitT func(key []byte, value string) error

// The mr library calls the reduce function once for each key
// generated by the map tasks, with a list of all the values created
// for that key by any map task.
type ReduceT func(string, []string, EmitT) error

// The mr library calls the map function for each line of input, which
// is passed in as an io.Reader.  The map function outputs its values
// by calling an emit function and passing it a KeyValue.
type MapT func(string, *bufio.Scanner, EmitT) error
