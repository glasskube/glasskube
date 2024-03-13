package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	"go.uber.org/multierr"
)

type EncodedKeyPair struct{ Key, Cert []byte }

func (kp *EncodedKeyPair) SaveTo(dir string) error {
	if err := os.MkdirAll(dir, 0740); err != nil {
		return err
	}
	var err error
	for name, data := range kp.AsMap() {
		err = multierr.Append(err, os.WriteFile(path.Join(dir, name), data, 0640))
	}
	return err
}

func (kp *EncodedKeyPair) AsMap() map[string][]byte {
	return map[string][]byte{
		"tls.key": kp.Key,
		"tls.crt": kp.Cert,
	}
}

type KeyPair struct {
	Key  *rsa.PrivateKey
	Cert *x509.Certificate
}

// Encoded returns the PEM encoded key pair
func (kp *KeyPair) Encoded() (*EncodedKeyPair, error) {
	var cb, kb bytes.Buffer

	if err := pem.Encode(&cb, &pem.Block{Type: "CERTIFICATE", Bytes: kp.Cert.Raw}); err != nil {
		return nil, fmt.Errorf("could not encode key: %w", err)
	}

	pkcs1 := x509.MarshalPKCS1PrivateKey(kp.Key)
	if err := pem.Encode(&kb, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: pkcs1}); err != nil {
		return nil, fmt.Errorf("could not encode cert: %w", err)
	}

	return &EncodedKeyPair{Key: kb.Bytes(), Cert: cb.Bytes()}, nil
}

// GenerateKeyPair generates a private key and certificate for a given certificate config.
// If ca is nil, a new self-signed CA keypair is generated.
func GenerateKeyPair(config *x509.Certificate, ca *KeyPair) (*KeyPair, error) {
	var kp KeyPair

	if pk, err := rsa.GenerateKey(rand.Reader, 4096); err != nil {
		return nil, err
	} else {
		kp.Key = pk
	}

	// When generating a CA key pair, parent must be the same as config and the generated private key is used
	parentCert := config
	parentKey := kp.Key
	if ca != nil {
		parentCert = ca.Cert
		parentKey = ca.Key
	}

	if certBytes, err := x509.CreateCertificate(
		rand.Reader, config, parentCert, &kp.Key.PublicKey, parentKey); err != nil {
		return nil, err
	} else if cert, err := x509.ParseCertificate(certBytes); err != nil {
		return nil, err
	} else {
		kp.Cert = cert
	}

	return &kp, nil
}
