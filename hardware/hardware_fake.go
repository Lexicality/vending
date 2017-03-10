// +build !linux,!arm

package hardware

import (
	"context"

	"github.com/lexicality/vending/vend"
	"github.com/op/go-logging"
)

type hardware struct {
	log *logging.Logger
}

func (hw *hardware) Setup() error {
	if hw.log != nil {
		hw.log.Info("Hello I'm not ARM!")
	}
	return nil
}

func (hw *hardware) Teardown() error {
	return nil
}

func (hw *hardware) Vend(ctx context.Context, location uint8) vend.Result {
	if hw.log != nil {
		hw.log.Infof("=== I AM VENDING ITEM #%d! ===", location)
	}
	return vend.ResultSuccess
}
