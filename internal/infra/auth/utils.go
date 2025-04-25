package auth

import (
	"bytes"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func ParsePublicKey(keyPem string) (*rsa.PublicKey, error) {
	return jwt.ParseRSAPublicKeyFromPEM([]byte(keyPem))
}

func ParsePrivateKey(keyPem string) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM([]byte(keyPem))
}

func ParseInlinePublicKey(key string) (*rsa.PublicKey, error) {
	segmentedKey := splitSubN(key, 64)
	keyPem := fmt.Sprintf("-----BEGIN PUBLIC KEY-----\n%s\n-----END PUBLIC KEY-----", strings.Join(segmentedKey, "\n"))
	return ParsePublicKey(keyPem)
}

func ParseInlinePrivateKey(key string) (*rsa.PrivateKey, error) {
	segmentedKey := splitSubN(key, 64)
	keyPem := fmt.Sprintf("-----BEGIN RSA PRIVATE KEY-----\n%s\n-----END RSA PRIVATE KEY-----", strings.Join(segmentedKey, "\n"))
	return ParsePrivateKey(keyPem)
}

func splitSubN(s string, n int) []string {
	sub := ""
	subs := []string{}

	runes := bytes.Runes([]byte(s))
	l := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			subs = append(subs, sub)
			sub = ""
		} else if (i + 1) == l {
			subs = append(subs, sub)
		}
	}

	return subs
}
