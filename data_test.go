package events

import (
	"encoding/json"
	"errors"
	"fmt"
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
		fields fields
		want   map[string]interface{}
	}{
		{
			fields: fields{
				event:   "event:happened",
				Payload: "some payload",
			},
			want: map[string]interface{}{
				"event":   "event:happened",
				"payload": "some payload",
				"message": "",
				"error":   "",
			},
		},
		{
			fields: fields{
				event: "event:happened",
				Error: errors.New("error"),
			},
			want: map[string]interface{}{
				"event":   "event:happened",
				"payload": "",
				"message": "",
				"error":   "error",
			},
		},
		{
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

	for i, tt := range tests {
		t.Run(fmt.Sprintf("marshalling(%d)", i), func(t *testing.T) {
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

func TestData_String(t *testing.T) {
	type fields struct {
		Payload interface{}
		Message string
		Error   error
		event   string
	}
	tests := []struct {
		fields fields
		want   map[string]interface{}
	}{
		{
			fields: fields{
				event:   "event:happened",
				Payload: "some payload",
			},
			want: map[string]interface{}{
				"event":   "event:happened",
				"payload": "some payload",
				"message": "",
				"error":   "",
			},
		},
		{
			fields: fields{
				event: "event:happened",
				Error: errors.New("error"),
			},
			want: map[string]interface{}{
				"event":   "event:happened",
				"payload": "",
				"message": "",
				"error":   "error",
			},
		},
		{
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

	for i, tt := range tests {
		t.Run(fmt.Sprintf("string(%d)", i), func(t *testing.T) {
			e := Data{
				Payload: tt.fields.Payload,
				Message: tt.fields.Message,
				Error:   tt.fields.Error,
				event:   tt.fields.event,
			}
			got := e.String()
			want, _ := json.Marshal(tt.want)

			if string(got) != string(want) {
				t.Errorf("Data.String(): expected = %v, got = %v", string(want), string(got))
			}
		})
	}
}
