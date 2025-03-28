package ticklens

import (
	"bytes"
	_ "embed"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	//go:embed TickLensProxy.json
	tickLensProxyJson []byte
	tickLensABI       abi.ABI

	//go:embed Events.json
	eventsJson []byte
	eventsABI  abi.ABI

	burnEvent abi.Event
)

func init() {
	builder := []struct {
		ABI  *abi.ABI
		data []byte
	}{
		{&tickLensABI, tickLensProxyJson},
		{&eventsABI, eventsJson},
	}

	for _, b := range builder {
		var err error
		*b.ABI, err = abi.JSON(bytes.NewReader(b.data))
		if err != nil {
			panic(err)
		}
	}
	burnEvent = eventsABI.Events["Burn"]
}
