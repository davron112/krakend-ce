package jwtvalidator

import (
	"api-gateway/v2/modules/lura/v2/logging"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

var (
	jwksCache   JWKS
	jwksCacheMu sync.Mutex
)

type JWKS struct {
	Keys []json.RawMessage `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func fetchJWKS(jwksURL string) (JWKS, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return JWKS{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return JWKS{}, errors.New("failed to fetch JWKS")
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return JWKS{}, err
	}
	return jwks, nil
}

func getJWK(keyID string, jwksURL string) (JWK, error) {
	jwksCacheMu.Lock()
	defer jwksCacheMu.Unlock()

	if len(jwksCache.Keys) == 0 {
		newJWKS, err := fetchJWKS(jwksURL)
		if err != nil {
			return JWK{}, err
		}
		jwksCache = newJWKS
	}

	for _, key := range jwksCache.Keys {
		var jwk JWK
		if err := json.Unmarshal(key, &jwk); err != nil {
			continue
		}
		if jwk.Kid == keyID {
			return jwk, nil
		}
	}

	// If kid not found, return the first JWK from the array
	if len(jwksCache.Keys) > 0 {
		var jwk JWK
		if err := json.Unmarshal(jwksCache.Keys[0], &jwk); err == nil {
			return jwk, nil
		}
	}

	return JWK{}, errors.New("key not found in JWKS and no keys available")
}

func validateJWT(cfg *Config, req *http.Request, l logging.Logger) (int, string) {
	var tokenString string
	headers := req.Header

	if headers.Get("X-User-Id") != "" || headers.Get("X-Session-Id") != "" {
		// Return 400 Bad Request
		return http.StatusBadRequest, "Invalid request headers"
	}

	tokenValue := headers.Get("Authorization")
	fmt.Println(tokenValue, "tokenValue")

	tokenParts := strings.Split(tokenValue, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return http.StatusUnauthorized, "Authorization header format must be 'Bearer {token}'"
	}
	tokenString = tokenParts[1]

	type CustomClaims struct {
		UserID    string   `json:"sub"`
		SessionID string   `json:"jti"`
		Roles     []string `json:"roles"`
		TokenType string   `json:"authenticationType"`
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if cfg.JwtPublicKey != "" {
			parsedKey, err := parsePublicKeyFromPEM(cfg.JwtPublicKey)
			if err != nil {
				return nil, fmt.Errorf("error parsing public key: %v", err)
			}
			return parsedKey, nil
		}

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return nil, errors.New("unexpected signing method")
		}

		keyID, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found in token")
		}

		jwk, err := getJWK(keyID, cfg.JWKSURL)
		if err != nil {
			return nil, err
		}

		return getKeyFromJWK(jwk)
	})

	if err != nil {
		fmt.Println(err)
		return http.StatusUnauthorized, "Unauthorized: Authentication is required and has failed or has not yet been provided."
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if cfg.TokenTypeField != "" && claims.TokenType != cfg.TokenTypeField {
			return http.StatusUnauthorized, "Unauthorized: Invalid token type."
		}
		if len(cfg.Roles) > 0 {
			roleMap := make(map[string]bool)
			for _, role := range cfg.Roles {
				roleMap[role] = true
			}

			// Check if the user has any of the required roles
			roleFound := false
			for _, role := range claims.Roles {
				if roleMap[role] {
					roleFound = true
					break
				}
			}

			if !roleFound {
				return http.StatusForbidden, "Forbidden: You do not have permission to access this resource with your current role."
			}
		}

		fmt.Println(claims, "claims")
		req.Header.Set("X-User-Id", claims.UserID)
		req.Header.Set("X-Session-Id", claims.SessionID)
		return 0, ""
	}

	return http.StatusUnauthorized, "Invalid token or claims"
}

func parsePublicKeyFromPEM(keyPEM string) (interface{}, error) {

	if !strings.Contains(keyPEM, "BEGIN") {
		keyPEM = "-----BEGIN PUBLIC KEY-----\n" + keyPEM + "\n-----END PUBLIC KEY-----"
	}

	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("unsupported key type")
	}
}

func getKeyFromJWK(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("invalid JWK N value: %v", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("invalid JWK E value: %v", err)
	}

	eInt := new(big.Int).SetBytes(eBytes).Int64()

	rsaPublicKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: int(eInt),
	}

	return rsaPublicKey, nil
}
