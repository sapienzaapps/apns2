package apns2_test

import (
	"encoding/json"
	"net/http"
	"testing"

	apns "github.com/sapienzaapps/apns2"
)

// Unit Tests

func TestResponseSent(t *testing.T) {
	if http.StatusOK != apns.StatusSent {
		t.Fatal("Expected:", apns.StatusSent, " found:", http.StatusOK)
	}
	if true != (&apns.Response{StatusCode: 200}).Sent() {
		t.Fatal("Expected:", (&apns.Response{StatusCode: 200}).Sent(), " found:", true)
	}
	if false != (&apns.Response{StatusCode: 400}).Sent() {
		t.Fatal("Expected:", (&apns.Response{StatusCode: 400}).Sent(), " found:", false)
	}
}

func TestIntTimestampParse(t *testing.T) {
	response := &apns.Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\":1458114061260}"
	_ = json.Unmarshal([]byte(payload), &response)
	if int64(1458114061260)/1000 != response.Timestamp.Unix() {
		t.Fatal("Expected:", response.Timestamp.Unix(), " found:", int64(1458114061260)/1000)
	}
}

func TestInvalidTimestampParse(t *testing.T) {
	response := &apns.Response{}
	payload := "{\"reason\":\"Unregistered\", \"timestamp\": \"2016-01-16 17:44:04 +1300\"}"
	err := json.Unmarshal([]byte(payload), &response)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
}
