package apns2_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/sapienzaapps/apns2"
)

func TestMarshalJSON(t *testing.T) {
	scenarios := []struct {
		in  interface{}
		out []byte
		err error
	}{
		{`{"a": "b"}`, []byte(`{"a": "b"}`), nil},
		{[]byte(`{"a": "b"}`), []byte(`{"a": "b"}`), nil},
		{struct {
			A string `json:"a"`
		}{"b"}, []byte(`{"a":"b"}`), nil},
	}

	notification := &apns2.Notification{}

	for _, scenario := range scenarios {
		notification.Payload = scenario.in
		payloadBytes, err := notification.MarshalJSON()

		if !bytes.Equal(scenario.out, payloadBytes) {
			t.Fatal("Expected:", payloadBytes, " found:", scenario.out)
		}
		if !errors.Is(scenario.err, err) {
			t.Fatal("Expected:", err, " found:", scenario.err)
		}
	}
}
