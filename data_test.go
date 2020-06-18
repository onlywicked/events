package events

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestData_MarshalJSON(t *testing.T) {
	type fields struct {
		Payload interface{}
		Message string
		Error   error
		event   string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{
			name: "should return only specified field",
			fields: fields{
				event:   "event:happened",
				Payload: "some payload",
			},
			want: map[string]interface{}{
				"event":   "event:happened",
				"payload": "some payload",
			},
		},
		{
			name: "should use error string",
			fields: fields{
				event: "event:happened",
				Error: errors.New("error"),
			},
			want: map[string]interface{}{
				"event": "event:happened",
				"error": "error",
			},
		},
		{
			name: "should return all fields",
			fields: fields{
				event:   "event:happened",
				Payload: "some payload",
				Message: "some message",
				Error:   errors.New("error"),
			},
			want: map[string]interface{}{
				"event":   "event:happened",
				"payload": "some payload",
				"message": "some message",
				"error":   "error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Data{
				Payload: tt.fields.Payload,
				Message: tt.fields.Message,
				Error:   tt.fields.Error,
				event:   tt.fields.event,
			}
			got, _ := e.MarshalJSON()
			want, _ := json.Marshal(tt.want)

			if string(got) != string(want) {
				t.Errorf("Data.MarshalJSON(): expected = %v, got = %v", string(want), string(got))
			}
		})
	}
}
