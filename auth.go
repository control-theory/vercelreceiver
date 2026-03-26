package vercelreceiver

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
)

const xVercelSignatureHeader = "x-vercel-signature"

// SignatureAlgorithm defines the supported signature algorithms
type SignatureAlgorithm string

const (
	SignatureAlgorithmSHA1   SignatureAlgorithm = "sha1"
	SignatureAlgorithmSHA256 SignatureAlgorithm = "sha256"
)

// verifySignature validates the x-vercel-signature header against the secret
func verifySignature(secret, bodyBytes []byte, signature string, algorithm string) bool {
	if len(secret) == 0 {
		// If no secret is configured, skip verification
		return true
	}

	if signature == "" {
		return false
	}

	var expectedMAC []byte
	switch algorithm {
	case "sha1":
		mac := hmac.New(sha1.New, secret)
		mac.Write(bodyBytes)
		expectedMAC = mac.Sum(nil)
	case "sha256":
		mac := hmac.New(sha256.New, secret)
		mac.Write(bodyBytes)
		expectedMAC = mac.Sum(nil)
	default:
		// Unknown algorithm, fail closed
		return false
	}

	expectedSignature := hex.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// verifyRequest validates the HTTP request signature
func verifyRequest(r *http.Request, secret string, bodyBytes []byte, algorithm string) error {
	if secret == "" {
		// No secret configured, skip verification
		return nil
	}

	signature := r.Header.Get(xVercelSignatureHeader)
	if signature == "" {
		return fmt.Errorf("missing %s header", xVercelSignatureHeader)
	}

	if !verifySignature([]byte(secret), bodyBytes, signature, algorithm) {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
