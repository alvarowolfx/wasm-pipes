package services

import (
	"github.com/apex/log"

	"github.com/alvarowolfx/wasm-pipes/engine"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type logPlugin struct {
	printFunc *wasmer.Function
	logger    *log.Entry
}

func init() {
	engine.RegisterPlugin("log", &logPlugin{})
}

func (p *logPlugin) Init(deps *engine.EngineDeps, wasmStore *wasmer.Store, vars *engine.ExecVars) error {
	p.logger = log.WithField("module", "wasm-log")
	p.printFunc = wasmer.NewFunction(
		wasmStore,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32), // a pointer for the message and length of message
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			msg := getString(vars.Instance, args[0], args[1])
			p.logger.Infof(msg)
			return []wasmer.Value{}, nil
		},
	)
	return nil
}

func (p *logPlugin) ImportedObjects() map[string]wasmer.IntoExtern {
	return map[string]wasmer.IntoExtern{
		"println": p.printFunc,
	}
}
