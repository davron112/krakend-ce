package krakend

import (
	"encoding/json"
	"errors"

	botdetector "api-gateway/v2/modules/krakend-botdetector/v2"
	"api-gateway/v2/modules/lura/v2/config"
)

// Namespace is the key used to store the bot detector config at the ExtraConfig struct
const Namespace = "github_com/davron112/krakend-botdetector"

// ErrNoConfig is returned when there is no config defined for the module
var ErrNoConfig = errors.New("no config defined for the module")

// ParseConfig extracts the module config from the ExtraConfig and returns a struct
// suitable for using the botdetector package
func ParseConfig(cfg config.ExtraConfig) (botdetector.Config, error) {
	res := botdetector.Config{}
	e, ok := cfg[Namespace]
	if !ok {
		return res, ErrNoConfig
	}
	b, err := json.Marshal(e)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &res)
	return res, err
}
