package builder

import (
	"github.com/primitivelab/goexchange"
	"testing"
)

func TestGetDepth(t *testing.T) {
	DefaultAPIBuilder.APIKey("")
	DefaultAPIBuilder.APISecretKey("")
	api := DefaultAPIBuilder.Build("mxc")
	t.Log(api.GetDepth(goexchange.NewSymbol("btc", "usdt"), 4, map[string]string{"type": "step0"}))
}

func TestGetUserBalance(t *testing.T) {

	DefaultAPIBuilder.APIKey("")
	DefaultAPIBuilder.APISecretKey("")
	DefaultAPIBuilder.Passphrase("")

	api := DefaultAPIBuilder.Build("okex")
	t.Log(api.GetUserBalance())
}
