package jwtvalidator

import (
	"encoding/json"
	"errors"

	"api-gateway/v2/modules/lura/v2/config"
)

const Namespace = "github_com/davron112/jwtvalidator"

type Config struct {
	JWKSURL              string   `json:"jwks_url"`
	JwtPublicKey         string   `json:"jwt_public_key"`
	AccessTokenHeaderKey string   `json:"access_token_header_key"`
	Roles                []string `json:"roles"`
	TokenTypeField       string   `json:"token_type_field"` // default value is "access-token"
}

var ErrNoConfig = errors.New("no JWT validator config")

func ParseConfig(extraConfig config.ExtraConfig) (*Config, error) {
	res := Config{}
	e, ok := extraConfig[Namespace]
	if !ok {
		return nil, ErrNoConfig
	}
	b, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &res)
	return &res, err
}
