package certificates

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"time"
)

var (
	// DefaultValidity is 365 days
	DefaultValidity = 365 * 24 * time.Hour
)

type Certificates struct{ Webhook, Ca *KeyPair }

func Generate(serviceName, serviceNamespace string, validity time.Duration) (*Certificates, error) {
	var certificates Certificates

	notBefore := time.Now()
	notAfter := notBefore.Add(validity)
	caCertificateConfig := &x509.Certificate{
		SerialNumber:          big.NewInt(0),
		Subject:               pkix.Name{Organization: []string{"glasskube.dev"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	if kp, err := GenerateKeyPair(caCertificateConfig, nil); err != nil {
		return nil, fmt.Errorf("could not create CA key pair: %w", err)
	} else {
		certificates.Ca = kp
	}

	serviceNameWithNamespace := fmt.Sprintf("%v.%v", serviceName, serviceNamespace)
	commonName := fmt.Sprintf("%v.svc", serviceNameWithNamespace)
	webhookCertificateConfig := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: commonName},
		DNSNames:     []string{serviceName, serviceNameWithNamespace, commonName},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	if kp, err := GenerateKeyPair(webhookCertificateConfig, certificates.Ca); err != nil {
		return nil, fmt.Errorf("could not create webhook key pair: %w", err)
	} else {
		certificates.Webhook = kp
	}

	return &certificates, nil
}
