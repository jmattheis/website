package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/rs/zerolog/log"
	"math/big"
	"time"
)

func NewUntrustedCert() *tls.Certificate {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		log.Fatal().Err(err).Msg("untrusted cert")
	}

	now := time.Now().UTC()

	subject := pkix.Name{
		CommonName:   "localhost",
		Organization: []string{"jmattheis.de"},
	}

	template := &x509.Certificate{
		SerialNumber:          big.NewInt(1337),
		Subject:               subject,
		DNSNames:              []string{"localhost"},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(time.Hour * 24 * 7),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  false,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)

	if err != nil {
		log.Fatal().Err(err).Msg("untrusted cert")
	}

	var certPem, keyPem bytes.Buffer
	if err2 := pem.Encode(&certPem, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err2 != nil {
		return nil
	}

	if err2 := pem.Encode(
		&keyPem,
		&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}); err2 != nil {
		return nil
	}

	c, err := tls.X509KeyPair(certPem.Bytes(), keyPem.Bytes())

	return &c
}
