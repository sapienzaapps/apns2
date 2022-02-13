package token_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sapienzaapps/apns2/token"
)

// AuthToken

func TestValidTokenFromP8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-valid.p8")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestValidTokenFromP8Bytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/authkey-valid.p8")
	_, err := token.AuthKeyFromBytes(bytes)
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
}

func TestNoSuchFileP8File(t *testing.T) {
	token, err := token.AuthKeyFromFile("")
	if errors.New("open : no such file or directory").Error() != err.Error() {
		t.Fatal("Expected:", err.Error(), " found:", errors.New("open : no such file or directory").Error())
	}
	if token != nil {
		t.Fatal("token expected nil, found:", token)
	}
}

func TestInvalidP8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-invalid.p8")
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
}

func TestInvalidPKCS8P8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-invalid-pkcs8.p8")
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
}

func TestInvalidECDSAP8File(t *testing.T) {
	_, err := token.AuthKeyFromFile("_fixtures/authkey-invalid-ecdsa.p8")
	if err == nil {
		t.Fatal("Expected error, found nil")
	}
}

// Expiry & Generation

func TestExpired(t *testing.T) {
	token := &token.Token{}
	if !token.Expired() {
		t.Fatal("Expected token expired")
	}
}

func TestNotExpired(t *testing.T) {
	token := &token.Token{
		IssuedAt: time.Now().Unix(),
	}
	if token.Expired() {
		t.Fatal("Expected token NOT expired")
	}
}

func TestExpiresBeforeAnHour(t *testing.T) {
	token := &token.Token{
		IssuedAt: time.Now().Add(-50 * time.Minute).Unix(),
	}
	if !token.Expired() {
		t.Fatal("Expected token expired")
	}
}

func TestGenerateIfExpired(t *testing.T) {
	authKey, _ := token.AuthKeyFromFile("_fixtures/authkey-valid.p8")
	token := &token.Token{
		AuthKey: authKey,
	}
	token.GenerateIfExpired()
	if time.Now().Unix() != token.IssuedAt {
		t.Fatal("Expected:", token.IssuedAt, " found:", time.Now().Unix())
	}
}

func TestGenerateWithNoAuthKey(t *testing.T) {
	token := &token.Token{}
	bool, err := token.Generate()
	if bool {
		t.Fatal("Expected generate false")
	} else if err == nil {
		t.Fatal("Expected error, found nil")
	}
}

func TestGenerateWithInvalidAuthKey(t *testing.T) {
	pubkeyCurve := elliptic.P521()
	privatekey, _ := ecdsa.GenerateKey(pubkeyCurve, rand.Reader)
	token := &token.Token{
		AuthKey: privatekey,
	}
	bool, err := token.Generate()
	if bool {
		t.Fatal("Expected generate false")
	} else if err == nil {
		t.Fatal("Expected error, found nil")
	}
}
