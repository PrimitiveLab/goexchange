package builder

import (
	"github.com/primitivelab/goexchange"
	"testing"
)

func TestGetDepth(t *testing.T) {

	api := DefaultAPIBuilder.Build("huobi")
	t.Log(api.GetDepth(goexchange.NewSymbol("btc", "usdt"), 4, map[string]string{"type":"step0"}))
}


