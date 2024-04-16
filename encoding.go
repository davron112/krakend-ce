package krakend

import (
	rss "api-gateway/v2/modules/krakend-rss/v2"
	xml "api-gateway/v2/modules/krakend-xml/v2"
	ginxml "api-gateway/v2/modules/krakend-xml/v2/gin"
	"api-gateway/v2/modules/lura/v2/router/gin"
)

// RegisterEncoders registers all the available encoders
func RegisterEncoders() {
	xml.Register()
	rss.Register()

	gin.RegisterRender(xml.Name, ginxml.Render)
}
