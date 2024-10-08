package jwtvalidator

import (
	"api-gateway/v2/modules/lura/v2/config"
	"api-gateway/v2/modules/lura/v2/logging"
	"api-gateway/v2/modules/lura/v2/proxy"
	krakendgin "api-gateway/v2/modules/lura/v2/router/gin"
	"github.com/gin-gonic/gin"
)

const logPrefix = "[SERVICE: Gin][JWTValidator]"

func Register(cfg config.ServiceConfig, l logging.Logger, engine *gin.Engine) {
	validatorCfg, err := ParseConfig(cfg.ExtraConfig)
	if err == ErrNoConfig {
		return
	}
	if err != nil {
		l.Warning(logPrefix, err.Error())
		return
	}

	l.Debug(logPrefix, "The JWT validator has been registered successfully")
	engine.Use(middleware(validatorCfg, l))
}

func New(hf krakendgin.HandlerFactory, l logging.Logger) krakendgin.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(cfg, p)
		logPrefix := "[ENDPOINT: " + cfg.Endpoint + "][JWTValidator]"

		validatorCfg, err := ParseConfig(cfg.ExtraConfig)
		if err == ErrNoConfig {
			return next
		}
		if err != nil {
			l.Warning(logPrefix, err.Error())
			return next
		}

		l.Debug(logPrefix, "The JWT validator has been registered successfully")
		return handler(validatorCfg, next, l)
	}
}
func middleware(cfg *Config, l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode, errMsg := validateJWT(cfg, c.Request, l)
		if statusCode != 0 {
			c.AbortWithStatusJSON(statusCode, gin.H{
				"error": gin.H{
					"msg": errMsg,
				},
			})
			return
		}

		c.Next()
	}
}

func handler(cfg *Config, next gin.HandlerFunc, l logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode, errMsg := validateJWT(cfg, c.Request, l)
		if statusCode != 0 {
			c.AbortWithStatusJSON(statusCode, gin.H{
				"error": gin.H{
					"msg": errMsg,
				},
			})
			return
		}

		next(c)
	}
}
