package core

import (
	"github.com/dapperlabs/bamboo-node/language/runtime"
	"github.com/dapperlabs/bamboo-node/pkg/types"

	"github.com/dapperlabs/bamboo-node/internal/emulator/state"
)

// Computer provides an interface to execute scripts against the world state.
type Computer struct {
	runtime runtime.Runtime
}

// NewComputer returns a new computer instance.
func NewComputer(runtime runtime.Runtime) *Computer {
	return &Computer{
		runtime: runtime,
	}
}

type runtimeInterface struct {
	getValue      func(controller, owner, key []byte) (value []byte, err error)
	setValue      func(controller, owner, key, value []byte) (err error)
	createAccount func(publicKey, code []byte) (id []byte, err error)
}

func (i *runtimeInterface) GetValue(controller, owner, key []byte) ([]byte, error) {
	return i.getValue(controller, owner, key)
}

func (i *runtimeInterface) SetValue(controller, owner, key, value []byte) error {
	return i.setValue(controller, owner, key, value)
}

func (i *runtimeInterface) CreateAccount(publicKey, code []byte) (id []byte, err error) {
	return i.createAccount(publicKey, code)
}

// ExecuteTransaction executes a transaction against the current world state.
func (c *Computer) ExecuteTransaction(
	tx *types.SignedTransaction,
	registers *state.RegistersView,
) (err error) {
	// TODO: deduct gas cost from transaction signer's account
	_, err = c.ExecuteScript(tx.Script, registers)
	return err
}

// ExecuteScript executes a script against the current world state.
func (c *Computer) ExecuteScript(
	script []byte,
	registers *state.RegistersView,
) (result interface{}, err error) {
	runtimeInterface := &runtimeInterface{
		getValue: func(controller, owner, key []byte) ([]byte, error) {
			v, _ := registers.Get(controller, owner, key)
			return v, nil
		},
		setValue: func(controller, owner, key, value []byte) error {
			registers.Set(controller, owner, key, value)
			return nil
		},
		createAccount: func(publicKey, code []byte) (id []byte, err error) {
			accountID := registers.CreateAccount(publicKey, code)
			return accountID.Bytes(), nil
		},
	}

	return c.runtime.ExecuteScript(script, runtimeInterface)
}
