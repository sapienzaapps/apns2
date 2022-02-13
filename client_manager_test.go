package apns2_test

import (
	"bytes"
	"crypto/tls"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/sapienzaapps/apns2"
	"github.com/sapienzaapps/apns2/certificate"
)

func TestNewClientManager(t *testing.T) {
	manager := apns2.NewClientManager()
	if manager.MaxSize != 64 {
		t.Fatal("Expected:", 64, " found:", manager.MaxSize)
	}
	if manager.MaxAge != 10*time.Minute {
		t.Fatal("Expected:", 10*time.Minute, " found:", manager.MaxAge)
	}
}

func TestClientManagerGetWithoutNew(t *testing.T) {
	manager := apns2.ClientManager{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	c1 := manager.Get(mockCert())
	c2 := manager.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	if c1 == nil {
		t.Fatal("c1 expected not nil, found nil")
	}
	if v1.Pointer() != v2.Pointer() {
		t.Fatal("Expected:", v2.Pointer(), " found:", v1.Pointer())
	}
	if 1 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 1)
	}
}

func TestClientManagerAddWithoutNew(t *testing.T) {
	wg := sync.WaitGroup{}

	manager := apns2.ClientManager{
		MaxSize: 1,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			manager.Add(apns2.NewClient(mockCert()))
			if 1 != manager.Len() {
				t.Error("Expected:", manager.Len(), " found:", 1)
				t.Fail()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestClientManagerLenWithoutNew(t *testing.T) {
	manager := apns2.ClientManager{
		MaxSize: 32,
		MaxAge:  5 * time.Minute,
		Factory: apns2.NewClient,
	}

	if 0 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 0)
	}
}

func TestClientManagerGetDefaultOptions(t *testing.T) {
	manager := apns2.NewClientManager()
	c1 := manager.Get(mockCert())
	c2 := manager.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	if c1 == nil {
		t.Fatal("c1 expected not nil, found nil")
	}
	if v1.Pointer() != v2.Pointer() {
		t.Fatal("Expected:", v2.Pointer(), " found:", v1.Pointer())
	}
	if 1 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 1)
	}
}

func TestClientManagerGetNilClientFactory(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.Factory = func(certificate tls.Certificate) *apns2.Client {
		return nil
	}
	c1 := manager.Get(mockCert())
	c2 := manager.Get(mockCert())
	if c1 != nil {
		t.Fatal("c1 expected nil, found:", c1)
	}
	if c2 != nil {
		t.Fatal("c2 expected nil, found:", c2)
	}
	if 0 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 0)
	}
}

func TestClientManagerGetMaxAgeExpiration(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.MaxAge = time.Nanosecond
	c1 := manager.Get(mockCert())
	time.Sleep(time.Microsecond)
	c2 := manager.Get(mockCert())
	v1 := reflect.ValueOf(c1)
	v2 := reflect.ValueOf(c2)
	if c1 == nil {
		t.Fatal("c1 expected not nil, found nil")
	}
	if v1.Pointer() == v2.Pointer() {
		t.Fatal("v1.Pointer() != v2.Pointer()")
	}
	if 1 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 1)
	}
}

func TestClientManagerGetMaxAgeExpirationWithNilFactory(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.Factory = func(certificate tls.Certificate) *apns2.Client {
		return nil
	}
	manager.MaxAge = time.Nanosecond
	manager.Add(apns2.NewClient(mockCert()))
	c1 := manager.Get(mockCert())
	time.Sleep(time.Microsecond)
	c2 := manager.Get(mockCert())
	if c1 != nil {
		t.Fatal("c1 expected nil, found:", c1)
	}
	if c2 != nil {
		t.Fatal("c2 expected nil, found:", c2)
	}
	if 1 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 1)
	}
}

func TestClientManagerGetMaxSizeExceeded(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.MaxSize = 1
	cert1 := mockCert()
	_ = manager.Get(cert1)
	cert2, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid.p12", "")
	_ = manager.Get(cert2)
	cert3, _ := certificate.FromP12File("certificate/_fixtures/certificate-valid-encrypted.p12", "password")
	c := manager.Get(cert3)
	if !bytes.Equal(cert3.Certificate[0], c.Certificate.Certificate[0]) {
		t.Fatal("Certificate not valid after decrypting")
	}
	if 1 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 1)
	}
}

func TestClientManagerAdd(t *testing.T) {
	fn := func(certificate tls.Certificate) *apns2.Client {
		t.Fatal("factory should not have been called")
		return nil
	}

	manager := apns2.NewClientManager()
	manager.Factory = fn
	manager.Add(apns2.NewClient(mockCert()))
	manager.Get(mockCert())
}

func TestClientManagerAddTwice(t *testing.T) {
	manager := apns2.NewClientManager()
	manager.Add(apns2.NewClient(mockCert()))
	manager.Add(apns2.NewClient(mockCert()))
	if 1 != manager.Len() {
		t.Fatal("Expected:", manager.Len(), " found:", 1)
	}
}
