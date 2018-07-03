package arc

import (
	"landzero.net/x/container/sack"
	"net/http"
)

// Context request context
type Context struct {
	sack.Sack

	Arc *Arc
	Req *Request
	Res ResponseWriter

	hIndex int // handler index
}

// Next execute next handler
func (c *Context) Next() {
	if c.hIndex >= len(c.Arc.handlers) {
		http.NotFound(c.Res, c.Req.Request)
		return
	}
	// save handler index
	i := c.hIndex
	// increase handler index
	c.hIndex++
	// call handler
	c.Arc.handlers[i](c)
}
