package lua

import (
	"api-gateway/v2/modules/binder"
)

func RegisterNil(b *binder.Binder) {
	tab := b.Table("luaNil")
	tab.Static("new", func(c *binder.Context) error {
		c.Push().Data(nil, "luaNil")
		return nil
	})
}
