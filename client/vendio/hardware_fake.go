// +build !linux,!arm

package vendio

import (
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

func (hw *hardware) Vend(location uint8) error {
	if hw.log != nil {
		hw.log.Infof("==== I AM VENDING ITEM #%d! ====", location)
	}
	return nil
}