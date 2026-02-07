package sso

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: Move this JWT validation code out into a separate package.

const wellKnownEndpoint = "https://login.eveonline.com/.well-known/oauth-authorization-server"
const cacheDuration = time.Minute * 15

type JWKCache struct {
	endpoint      string
	lastFetchedAt time.Time
	jwk           *JWK
	mutex         sync.Mutex
}

type JWK struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"` // base64url (RSA modulus)
	E   string `json:"e"` // base64url (RSA exponent)
}

func (s *SSO) validateToken(tokenString string) (jwt.Claims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}

	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return nil, err
	}

	s.jwkCache.mutex.Lock()
	cachedJWK := s.jwkCache.jwk
	lastFetchedAt := s.jwkCache.lastFetchedAt
	s.jwkCache.mutex.Unlock()

	if cachedJWK == nil || time.Since(lastFetchedAt) > cacheDuration || cachedJWK.Kid != header.Kid {
		err = s.getJWK(header.Kid)
		if err != nil {
			err = s.getJWK(header.Kid)
			if err != nil {
				return nil, err
			}
		}
	}

	s.jwkCache.mutex.Lock()
	jwk := s.jwkCache.jwk
	s.jwkCache.mutex.Unlock()

	publicKey, err := rsaPublicKeyFromJWK(*jwk)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	expirationTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}
	if time.Now().After(expirationTime.Time) {
		return nil, errors.New("token expired")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, err
	}
	if issuer != "https://login.eveonline.com" {
		return nil, errors.New("invalid issuer")
	}

	audience, err := token.Claims.GetAudience()
	if err != nil {
		return nil, err
	}
	if !slices.Contains(audience, "EVE Online") || !slices.Contains(audience, s.credentials.ClientID) {
		return nil, errors.New("invalid audience")
	}

	return token.Claims, nil
}

func (s *SSO) getJWK(kid string) error {
	err := s.getJWKSEndpoint()
	if err != nil {
		return err
	}

	response, err := s.client.Get(s.jwkCache.endpoint)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var jwks struct {
		Keys []JWK `json:"keys"`
	}
	err = json.Unmarshal(data, &jwks)
	if err != nil {
		return err
	}

	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			s.jwkCache.mutex.Lock()
			s.jwkCache.jwk = &jwk
			s.jwkCache.lastFetchedAt = time.Now()
			s.jwkCache.mutex.Unlock()

			return nil
		}
	}

	return errors.New("jwk not found for kid")
}

func (s *SSO) getJWKSEndpoint() error {
	s.jwkCache.mutex.Lock()
	endpoint := s.jwkCache.endpoint
	s.jwkCache.mutex.Unlock()

	if endpoint != "" {
		return nil
	}

	response, err := s.client.Get(wellKnownEndpoint)
	if err != nil {
		return errors.New("could not get JWKS endpoint")
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("unexpected status code: " + strconv.Itoa(response.StatusCode))
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.New("could not read JWKS endpoint")
	}

	var body struct {
		JwksURI string `json:"jwks_uri"`
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		return errors.New("could not read JWKS endpoint")
	}

	if body.JwksURI == "" {
		return errors.New("could not read JWKS endpoint")
	}

	s.jwkCache.mutex.Lock()
	s.jwkCache.endpoint = body.JwksURI
	s.jwkCache.mutex.Unlock()

	return nil
}

// Shamelessly copied from ChatGPT because cryptography gives me a headache...
func rsaPublicKeyFromJWK(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}
