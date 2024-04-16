package gin

import (
	"github.com/clbanning/mxj"
	"api-gateway/v2/modules/lura/v2/proxy"
	"github.com/gin-gonic/gin"
)

func init() {
	mxj.XMLEscapeChars(true)
}

// Render marshals the proxy response and passes the resulting xml to the response writer
func Render(c *gin.Context, response *proxy.Response) {
	status := c.Writer.Status()
	if response == nil {
		c.XML(status, nil)
		return
	}
	mv := mxj.Map(response.Data)
	c.Header("Content-Type", gin.MIMEXML)
	mv.XmlWriter(c.Writer)
}
