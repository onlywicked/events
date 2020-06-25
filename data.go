package events

import "encoding/json"

// Data represents an event data in a event lifecycle
type Data struct {
	event   string
	Payload interface{}
	// Optional
	Message string
	Error   error
}

// String implements the stringer interface for the Data
func (e Data) String() string {
	j, _ := e.MarshalJSON()
	return string(j)
}

// MarshalJSON implements the json marshaler interface for the Data
func (e Data) MarshalJSON() ([]byte, error) {
	eM := make(map[string]interface{})

	eM["event"] = e.event

	if e.Payload != nil {
		eM["payload"] = e.Payload
	} else {
		eM["payload"] = ""
	}

	eM["message"] = e.Message

	if e.Error != nil {
		eM["error"] = e.Error.Error()
	} else {
		eM["error"] = ""
	}

	return json.Marshal(eM)
}
