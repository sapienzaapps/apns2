package apns2_test

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/http2"

	apns "github.com/sapienzaapps/apns2"
	"github.com/sapienzaapps/apns2/certificate"
	"github.com/sapienzaapps/apns2/token"
)

// Mocks

func mockNotification() *apns.Notification {
	n := &apns.Notification{}
	n.DeviceToken = "11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"
	n.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	return n
}

func mockToken() *token.Token {
	pubkeyCurve := elliptic.P256()
	authKey, _ := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	return &token.Token{AuthKey: authKey}
}

func mockCert() tls.Certificate {
	return tls.Certificate{}
}

func mockClient(url string) *apns.Client {
	return &apns.Client{Host: url, HTTPClient: http.DefaultClient}
}

type mockTransport struct {
	*http2.Transport
	closed bool
}

func (c *mockTransport) CloseIdleConnections() {
	c.closed = true
}

// Unit Tests

func TestClientDefaultHost(t *testing.T) {
	client := apns.NewClient(mockCert())
	if "https://api.sandbox.push.apple.com" != client.Host {
		t.Fatal("Expected:", client.Host, " found:", "https://api.sandbox.push.apple.com")
	}
}

func TestTokenDefaultHost(t *testing.T) {
	client := apns.NewTokenClient(mockToken()).Development()
	if "https://api.sandbox.push.apple.com" != client.Host {
		t.Fatal("Expected:", client.Host, " found:", "https://api.sandbox.push.apple.com")
	}
}

func TestClientDevelopmentHost(t *testing.T) {
	client := apns.NewClient(mockCert()).Development()
	if "https://api.sandbox.push.apple.com" != client.Host {
		t.Fatal("Expected:", client.Host, " found:", "https://api.sandbox.push.apple.com")
	}
}

func TestTokenClientDevelopmentHost(t *testing.T) {
	client := apns.NewTokenClient(mockToken()).Development()
	if "https://api.sandbox.push.apple.com" != client.Host {
		t.Fatal("Expected:", client.Host, " found:", "https://api.sandbox.push.apple.com")
	}
}

func TestClientProductionHost(t *testing.T) {
	client := apns.NewClient(mockCert()).Production()
	if "https://api.push.apple.com" != client.Host {
		t.Fatal("Expected:", client.Host, " found:", "https://api.push.apple.com")
	}
}

func TestTokenClientProductionHost(t *testing.T) {
	client := apns.NewTokenClient(mockToken()).Production()
	if "https://api.push.apple.com" != client.Host {
		t.Fatal("Expected:", client.Host, " found:", "https://api.push.apple.com")
	}
}

func TestClientBadUrlError(t *testing.T) {
	n := mockNotification()
	res, err := mockClient("badurl://badurl.com").Push(n)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
	if res != nil {
		t.Fatal("res expected nil, found:", res)
	}
}

func TestClientBadTransportError(t *testing.T) {
	n := mockNotification()
	client := mockClient("badurl://badurl.com")
	client.HTTPClient.Transport = nil
	res, err := client.Push(n)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
	if res != nil {
		t.Fatal("res expected nil, found:", res)
	}
}

func TestClientBadDeviceToken(t *testing.T) {
	n := &apns.Notification{}
	n.DeviceToken = "DGw\aOoD+HwSroh#Ug]%xzd]"
	n.Payload = []byte(`{"aps":{"alert":"Hello!"}}`)
	res, err := mockClient("https://api.push.apple.com").Push(n)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
	if res != nil {
		t.Fatal("res expected nil, found:", res)
	}
}

func TestClientNameToCertificate(t *testing.T) {
	crt, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	client := apns.NewClient(crt)
	name := client.HTTPClient.Transport.(*http2.Transport).TLSClientConfig.Certificates
	if len(name) != 1 || len(name[0].Certificate) != 1 {
		t.Fatal("Expected 1 certificate")
	}

	certificate2 := tls.Certificate{}
	client2 := apns.NewClient(certificate2)
	name2 := client2.HTTPClient.Transport.(*http2.Transport).TLSClientConfig.Certificates
	if len(name2) != 1 || len(name2[0].Certificate) != 0 {
		t.Fatal("Expected zero certificates")
	}
}

func TestDialTLSTimeout(t *testing.T) {
	apns.TLSDialTimeout = 10 * time.Millisecond
	crt, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	client := apns.NewClient(crt)
	dialTLS := client.HTTPClient.Transport.(*http2.Transport).DialTLS
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	defer listener.Close()
	var e error
	if _, e = dialTLS("tcp", address, nil); e == nil {
		t.Fatal("Dial completed successfully")
	}
	if !strings.Contains(e.Error(), "timed out") && !strings.Contains(e.Error(), "context deadline exceeded") {
		t.Errorf("resulting error not a timeout: %s", e)
	}
}

// Functional Tests

func TestURL(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "POST" != r.Method {
			t.Fatal("Expected:", r.Method, " found:", "POST")
		}
		if fmt.Sprintf("/3/device/%s", n.DeviceToken) != r.URL.String() {
			t.Fatal("Expected:", r.URL.String(), " found:", fmt.Sprintf("/3/device/%s", n.DeviceToken))
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestDefaultHeaders(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "application/json; charset=utf-8" != r.Header.Get("Content-Type") {
			t.Fatal("Expected:", r.Header.Get("Content-Type"), " found:", "application/json; charset=utf-8")
		}
		if "" != r.Header.Get("apns-id") {
			t.Fatal("Expected:", r.Header.Get("apns-id"), " found:", "")
		}
		if "" != r.Header.Get("apns-collapse-id") {
			t.Fatal("Expected:", r.Header.Get("apns-collapse-id"), " found:", "")
		}
		if "" != r.Header.Get("apns-priority") {
			t.Fatal("Expected:", r.Header.Get("apns-priority"), " found:", "")
		}
		if "" != r.Header.Get("apns-topic") {
			t.Fatal("Expected:", r.Header.Get("apns-topic"), " found:", "")
		}
		if "" != r.Header.Get("apns-expiration") {
			t.Fatal("Expected:", r.Header.Get("apns-expiration"), " found:", "")
		}
		if "" != r.Header.Get("thread-id") {
			t.Fatal("Expected:", r.Header.Get("thread-id"), " found:", "")
		}
		if "alert" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "alert")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestClientPushWithContextWithTimeout(t *testing.T) {
	const timeout = time.Nanosecond
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	time.Sleep(timeout)
	res, err := mockClient(server.URL).PushWithContext(ctx, n)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
	if res != nil {
		t.Fatal("res expected nil, found:", res)
	}
	cancel()
}

func TestClientPushWithContext(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	res, err := mockClient(server.URL).PushWithContext(context.Background(), n)
	if err != nil {
		t.Fatal("err expected nil, found:", err)
	}
	if res.ApnsID != apnsID {
		t.Fatal("Expected:", apnsID, " found:", res.ApnsID)
	}
}

func TestHeaders(t *testing.T) {
	n := mockNotification()
	n.ApnsID = "84DB694F-464F-49BD-960A-D6DB028335C9"
	n.CollapseID = "game1.start.identifier"
	n.Topic = "com.testapp"
	n.Priority = 10
	n.Expiration = time.Now()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if n.ApnsID != r.Header.Get("apns-id") {
			t.Fatal("Expected:", r.Header.Get("apns-id"), " found:", n.ApnsID)
		}
		if n.CollapseID != r.Header.Get("apns-collapse-id") {
			t.Fatal("Expected:", r.Header.Get("apns-collapse-id"), " found:", n.CollapseID)
		}
		if "10" != r.Header.Get("apns-priority") {
			t.Fatal("Expected:", r.Header.Get("apns-priority"), " found:", "10")
		}
		if n.Topic != r.Header.Get("apns-topic") {
			t.Fatal("Expected:", r.Header.Get("apns-topic"), " found:", n.Topic)
		}
		if fmt.Sprintf("%v", n.Expiration.Unix()) != r.Header.Get("apns-expiration") {
			t.Fatal("Expected:", r.Header.Get("apns-expiration"), " found:", fmt.Sprintf("%v", n.Expiration.Unix()))
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPushTypeAlertHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeAlert
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "alert" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "alert")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPushTypeBackgroundHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeBackground
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "background" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "background")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPushTypeVOIPHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeVOIP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "voip" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "voip")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPushTypeComplicationHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeComplication
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "complication" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "complication")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPushTypeFileProviderHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeFileProvider
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "fileprovider" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "fileprovider")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPushTypeMDMHeader(t *testing.T) {
	n := mockNotification()
	n.PushType = apns.PushTypeMDM
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "mdm" != r.Header.Get("apns-push-type") {
			t.Fatal("Expected:", r.Header.Get("apns-push-type"), " found:", "mdm")
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestAuthorizationHeader(t *testing.T) {
	n := mockNotification()
	token := mockToken()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "application/json; charset=utf-8" != r.Header.Get("Content-Type") {
			t.Fatal("Expected:", r.Header.Get("Content-Type"), " found:", "application/json; charset=utf-8")
		}
		if fmt.Sprintf("bearer %v", token.Bearer) != r.Header.Get("authorization") {
			t.Fatal("Expected:", r.Header.Get("authorization"), " found:", fmt.Sprintf("bearer %v", token.Bearer))
		}
	}))
	defer server.Close()

	client := mockClient(server.URL)
	client.Token = token
	_, err := client.Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestPayload(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal("Expected no error, found:", err)
		}
		if !bytes.Equal(n.Payload.([]byte), body) {
			t.Fatal("Expected:", body, " found:", n.Payload)
		}
	}))
	defer server.Close()
	_, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestBadPayload(t *testing.T) {
	n := mockNotification()
	n.Payload = func() {}
	_, err := mockClient("").Push(n)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
}

func Test200SuccessResponse(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if http.StatusOK != res.StatusCode {
		t.Fatal("Expected:", res.StatusCode, " found:", http.StatusOK)
	}
	if apnsID != res.ApnsID {
		t.Fatal("Expected:", res.ApnsID, " found:", apnsID)
	}
	if true != res.Sent() {
		t.Fatal("Expected:", res.Sent(), " found:", true)
	}
}

func Test400BadRequestPayloadEmptyResponse(t *testing.T) {
	n := mockNotification()
	var apnsID = "02ABC856-EF8D-4E49-8F15-7B8A61D978D6"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("{\"reason\":\"PayloadEmpty\"}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if 400 != res.StatusCode {
		t.Fatal("Expected:", res.StatusCode, " found:", 400)
	}
	if apnsID != res.ApnsID {
		t.Fatal("Expected:", res.ApnsID, " found:", apnsID)
	}
	if apns.ReasonPayloadEmpty != res.Reason {
		t.Fatal("Expected:", res.Reason, " found:", apns.ReasonPayloadEmpty)
	}
	if false != res.Sent() {
		t.Fatal("Expected:", res.Sent(), " found:", false)
	}
}

func Test410UnregisteredResponse(t *testing.T) {
	n := mockNotification()
	var apnsID = "9F595474-356C-485E-B67F-9870BAE68702"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("apns-id", apnsID)
		w.WriteHeader(http.StatusGone)
		_, _ = w.Write([]byte("{\"reason\":\"Unregistered\", \"timestamp\": 1458114061260 }"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if 410 != res.StatusCode {
		t.Fatal("Expected:", res.StatusCode, " found:", 410)
	}
	if apnsID != res.ApnsID {
		t.Fatal("Expected:", res.ApnsID, " found:", apnsID)
	}
	if apns.ReasonUnregistered != res.Reason {
		t.Fatal("Expected:", res.Reason, " found:", apns.ReasonUnregistered)
	}
	if int64(1458114061260)/1000 != res.Timestamp.Unix() {
		t.Fatal("Expected:", res.Timestamp.Unix(), " found:", int64(1458114061260)/1000)
	}
	if false != res.Sent() {
		t.Fatal("Expected:", res.Sent(), " found:", false)
	}
}

func TestMalformedJSONResponse(t *testing.T) {
	n := mockNotification()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte("{{MalformedJSON}}"))
	}))
	defer server.Close()
	res, err := mockClient(server.URL).Push(n)
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
	if false != res.Sent() {
		t.Fatal("Expected:", res.Sent(), " found:", false)
	}
}

func TestCloseIdleConnections(t *testing.T) {
	transport := &mockTransport{}

	client := mockClient("")
	client.HTTPClient.Transport = transport

	if false != transport.closed {
		t.Fatal("Expected:", transport.closed, " found:", false)
	}
	client.CloseIdleConnections()
	if true != transport.closed {
		t.Fatal("Expected:", transport.closed, " found:", true)
	}
}
