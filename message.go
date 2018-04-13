package hub

import (
	"log"
	"time"
)

type (
	// Fields is a [key]value storage for Messages.
	Fields map[string]interface{}

	// Message represent some message/event passed into the hub
	// It also contain some helper functions to convert the fields to primitive types.
	Message struct {
		Name   string
		Msg    []byte
		Fields Fields
	}
)

// Topic return the message topic used when the message was sended.
func (m *Message) Topic() string {
	return m.Name
}

// Int return an int field of an Message.
func (m *Message) Int(key string) int {
	v, ok := m.Fields[key].(int)
	if !ok {
		log.Fatalf("Message %#v didn`t have the int field %s", m, key)
	}
	return v
}

// Int64 return an int64 field of an Message.
func (m *Message) Int64(key string) int64 {
	v, ok := m.Fields[key].(int64)
	if !ok {
		log.Fatalf("Message %#v didn`t have the int64 field %s", m, key)
	}
	return v
}

// Int32 return an int32 field of an Message.
func (m *Message) Int32(key string) int32 {
	v, ok := m.Fields[key].(int32)
	if !ok {
		log.Fatalf("Message %#v didn`t have the int32 field %s", m, key)
	}
	return v
}

// Float64 return a float64 field of an Message.
func (m *Message) Float64(key string) float64 {
	v, ok := m.Fields[key].(float64)
	if !ok {
		log.Fatalf("Message %#v didn`t have the float64 field %s", m, key)
	}
	return v
}

// String return a string field of an Message.
func (m *Message) String(key string) string {
	v, ok := m.Fields[key].(string)
	if !ok {
		log.Fatalf("Message %#v didn`t have the string field %s", m, key)
	}
	return v
}

// Duration return a time.Duration field of an Message.
func (m *Message) Duration(key string) time.Duration {
	v, ok := m.Fields[key].(time.Duration)
	if !ok {
		log.Fatalf("Message %#v didn`t have the time.Duration field %s", m, key)
	}
	return v
}
