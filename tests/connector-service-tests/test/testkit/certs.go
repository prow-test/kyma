package testkit

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

const (
	rsaKeySize = 2048
)

// Create Key generates rsa.PrivateKey
func CreateKey(t *testing.T) *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	require.NoError(t, err)

	return key
}

// CreateCsr creates CSR request
func CreateCsr(t *testing.T, certInfo CertInfo, keys *rsa.PrivateKey) []byte {
	subjectInfo := extractSubject(certInfo.Subject)

	subject := pkix.Name{
		CommonName:         subjectInfo["CN"],
		Country:            []string{subjectInfo["C"]},
		Organization:       []string{subjectInfo["O"]},
		OrganizationalUnit: []string{subjectInfo["OU"]},
		Locality:           []string{subjectInfo["L"]},
		Province:           []string{subjectInfo["ST"]},
	}

	var csrTemplate = x509.CertificateRequest{
		Subject: subject,
	}

	// step: generate the csr request
	csrCertificate, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keys)
	require.NoError(t, err)

	csr := pem.EncodeToMemory(&pem.Block{
		Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
	})

	return csr
}

// CrtResponseToPemBytes decodes certificate form CrtResponse and return pemBlock.Bytes
func CrtResponseToPemBytes(t *testing.T, certResponse *CrtResponse) []byte {
	crtBytes, err := base64.StdEncoding.DecodeString(certResponse.Crt)
	require.NoError(t, err)

	pemBlock, _ := pem.Decode(crtBytes)
	require.NotNil(t, pemBlock)

	return pemBlock.Bytes
}

// DecodeAndParseCert decodes base64 encoded certificate and parses it
func DecodeAndParseCert(t *testing.T, crtResponse *CrtResponse) *x509.Certificate {
	certBytes := CrtResponseToPemBytes(t, crtResponse)

	certificate, err := x509.ParseCertificate(certBytes)
	require.NoError(t, err)

	return certificate
}

// CheckIfSubjectEquals verifies that specified subject is equal to this in certificate
func CheckIfSubjectEquals(t *testing.T, expectedSubject string, certificate *x509.Certificate) {
	subjectInfo := extractSubject(expectedSubject)
	actualSubject := certificate.Subject

	require.Equal(t, subjectInfo["CN"], actualSubject.CommonName)
	require.Equal(t, []string{subjectInfo["C"]}, actualSubject.Country)
	require.Equal(t, []string{subjectInfo["O"]}, actualSubject.Organization)
	require.Equal(t, []string{subjectInfo["OU"]}, actualSubject.OrganizationalUnit)
	require.Equal(t, []string{subjectInfo["L"]}, actualSubject.Locality)
	require.Equal(t, []string{subjectInfo["ST"]}, actualSubject.Province)
}

func EncodeBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func extractSubject(subject string) map[string]string {
	result := map[string]string{}

	segments := strings.Split(subject, ",")

	for _, segment := range segments {
		parts := strings.Split(segment, "=")
		result[parts[0]] = parts[1]
	}

	return result
}
