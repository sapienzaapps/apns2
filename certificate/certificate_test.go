package certificate_test

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/sapienzaapps/apns2/certificate"
)

// PKCS#12

func TestValidCertificateFromP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid.p12", "")
	if err != nil {
		t.Fatal("err expected nil, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestValidCertificateFromP12Bytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.p12")
	cer, err := certificate.FromP12Bytes(bytes, "")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestEncryptedValidCertificateFromP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid-encrypted.p12", "password")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestNoSuchFileP12File(t *testing.T) {
	cer, err := certificate.FromP12File("", "")
	if errors.New("open : no such file or directory").Error() != err.Error() {
		t.Fatal("Expected:", err.Error(), " found:", errors.New("open : no such file or directory").Error())
	}
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
}

func TestBadPasswordP12File(t *testing.T) {
	cer, err := certificate.FromP12File("_fixtures/certificate-valid-encrypted.p12", "")
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
	if errors.New("pkcs12: decryption password incorrect").Error() != err.Error() {
		t.Fatal("Expected:", err.Error(), " found:", errors.New("pkcs12: decryption password incorrect").Error())
	}
}

// PEM

func TestValidCertificateFromPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid.pem", "")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestValidCertificateFromPemBytes(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid.pem")
	cer, err := certificate.FromPemBytes(bytes, "")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestValidCertificateFromPemFileWithPKCS8PrivateKey(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-pkcs8.pem", "")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestValidCertificateFromPemBytesWithPKCS8PrivateKey(t *testing.T) {
	bytes, _ := ioutil.ReadFile("_fixtures/certificate-valid-pkcs8.pem")
	cer, err := certificate.FromPemBytes(bytes, "")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestEncryptedValidCertificateFromPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-encrypted.pem", "password")
	if err != nil {
		t.Fatal("Expected no error, found:", err)
	}
	if certEqual(tls.Certificate{}, cer) {
		t.Fatal("tls.Certificate{} == cer")
	}
}

func TestNoSuchFilePemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("", "")
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
	if errors.New("open : no such file or directory").Error() != err.Error() {
		t.Fatal("Expected:", err.Error(), " found:", errors.New("open : no such file or directory").Error())
	}
}

func TestBadPasswordPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-valid-encrypted.pem", "badpassword")
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
	if !errors.Is(err, certificate.ErrFailedToDecryptKey) {
		t.Fatal("Expected:", err, " found:", certificate.ErrFailedToDecryptKey)
	}
}

func TestBadKeyPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-bad-key.pem", "")
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
	if !errors.Is(err, certificate.ErrFailedToParsePrivateKey) {
		t.Fatal("Expected:", err, " found:", certificate.ErrFailedToParsePrivateKey)
	}
}

func TestNoKeyPemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-no-key.pem", "")
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
	if !errors.Is(err, certificate.ErrNoPrivateKey) {
		t.Fatal("Expected:", err, " found:", certificate.ErrNoPrivateKey)
	}
}

func TestNoCertificatePemFile(t *testing.T) {
	cer, err := certificate.FromPemFile("_fixtures/certificate-no-certificate.pem", "")
	if !certEqual(tls.Certificate{}, cer) {
		t.Fatal("Expected:", cer, " found:", tls.Certificate{})
	}
	if !errors.Is(err, certificate.ErrNoCertificate) {
		t.Fatal("Expected:", err, " found:", certificate.ErrNoCertificate)
	}
}

func certEqual(cert tls.Certificate, buf tls.Certificate) bool {
	if (buf.Certificate == nil && cert.Certificate != nil) || (buf.Certificate != nil && cert.Certificate == nil) {
		return false
	}
	for _, b := range buf.Certificate {
		for _, c := range cert.Certificate {
			if !bytes.Equal(b, c) {
				return false
			}
		}
	}
	return true
}
