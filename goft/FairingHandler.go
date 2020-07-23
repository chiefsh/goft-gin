package goft

import (
	"github.com/gin-gonic/gin"
	"sync"
)

var fairingHandler *FairingHandler
var fairing_once sync.Once

func getFairingHandler() *FairingHandler {
	fairing_once.Do(func() {
		fairingHandler = &FairingHandler{}
	})
	return fairingHandler
}

type FairingHandler struct {
	fairings []Fairing
}

func (this *FairingHandler) AddFairing(f ...Fairing) {
	this.fairings = append(this.fairings, f...)
}
func (this *FairingHandler) before(ctx *gin.Context) {
	for _, f := range this.fairings {
		err := f.OnRequest(ctx)
		if err != nil {
			Throw(err.Error(), 400, ctx)
		}
	}
}
func (this *FairingHandler) after(ctx *gin.Context, ret interface{}) interface{} {
	var result = ret
	for _, f := range this.fairings {
		r, err := f.OnResponse(result)
		if err != nil {
			Throw(err.Error(), 400, ctx)
		}
		result = r
	}
	return result
}
func HandleFairing(responder Responder, ctx *gin.Context) interface{} {
	getFairingHandler().before(ctx)
	var ret interface{}
	if s1, ok := responder.(StringResponder); ok {
		ret = s1(ctx)
	}
	if s2, ok := responder.(JsonResponder); ok {
		ret = s2(ctx)
	}

	return getFairingHandler().after(ctx, ret)

}