package hub

type (
	// Fields is a [key]value storage for Messages values.
	Fields map[string]interface{}

	// Message represent some message/event passed into the hub
	// It also contain some helper functions to convert the fields to primitive types.
	Message struct {
		Name   string
		Body   []byte
		Fields Fields
	}
)

// Topic return the message topic used when the message was sended.
func (m *Message) Topic() string {
	return m.Name
}
