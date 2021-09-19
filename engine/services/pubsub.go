package services

import (
	"context"

	"github.com/alvarowolfx/wasm-pipes/engine"
	"github.com/apex/log"
	"github.com/wasmerio/wasmer-go/wasmer"
	"gocloud.dev/pubsub"
)

type pubsubPlugin struct {
	getMessageFunc *wasmer.Function
	sendFunc       *wasmer.Function
}

func init() {
	engine.RegisterPlugin("pubsub", &pubsubPlugin{})
}

func (p *pubsubPlugin) Init(deps *engine.EngineDeps, wasmStore *wasmer.Store, vars *engine.ExecVars) error {
	p.getMessageFunc = wasmer.NewFunction(
		wasmStore,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32), // a pointer to value response
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			setString(vars.Instance, args[0], string(vars.Message.Body))
			return []wasmer.Value{}, nil
		},
	)
	p.sendFunc = wasmer.NewFunction(
		wasmStore,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32), // a pointer for the message and length of message
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			msg := getString(vars.Instance, args[0], args[1])
			err := deps.Output.Send(context.Background(), &pubsub.Message{
				Body: []byte(msg),
			})
			if err != nil {
				log.Errorf("error sending message from wasm: %v", err)
			}
			return []wasmer.Value{}, nil
		},
	)
	return nil
}

func (p *pubsubPlugin) ImportedObjects() map[string]wasmer.IntoExtern {
	return map[string]wasmer.IntoExtern{
		"getMessage": p.getMessageFunc,
		"send":       p.sendFunc,
	}
}
