package services

import (
	"fmt"

	"github.com/alvarowolfx/wasm-pipes/engine"
	"github.com/eko/gocache/store"
	"github.com/wasmerio/wasmer-go/wasmer"
)

type cachePlugin struct {
	cacheGetKey *wasmer.Function
	cacheSetKey *wasmer.Function
}

func init() {
	engine.RegisterPlugin("cache", &cachePlugin{})
}

func (p *cachePlugin) Init(deps *engine.EngineDeps, wasmStore *wasmer.Store, vars *engine.ExecVars) error {
	p.cacheGetKey = wasmer.NewFunction(
		wasmStore,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32), // a pointer for the key and length of key, pointer to value response
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			key := getString(vars.Instance, args[0], args[1])
			value, cerr := deps.Cache.Get(key)
			if cerr != nil {
				setString(vars.Instance, args[2], "")
				return []wasmer.Value{}, nil
			}
			setString(vars.Instance, args[2], fmt.Sprintf("%s", value))
			return []wasmer.Value{}, nil
		},
	)
	p.cacheSetKey = wasmer.NewFunction(
		wasmStore,
		wasmer.NewFunctionType(
			wasmer.NewValueTypes(wasmer.I32, wasmer.I32, wasmer.I32, wasmer.I32), // a pointer for the key and length of key, pointer to value and length of value
			wasmer.NewValueTypes(),
		),
		func(args []wasmer.Value) ([]wasmer.Value, error) {
			key := getString(vars.Instance, args[0], args[1])
			value := getString(vars.Instance, args[2], args[3])
			_ = deps.Cache.Set(key, value, &store.Options{})
			return []wasmer.Value{}, nil
		},
	)
	return nil
}

func (p *cachePlugin) ImportedObjects() map[string]wasmer.IntoExtern {
	return map[string]wasmer.IntoExtern{
		"get": p.cacheGetKey,
		"set": p.cacheSetKey,
	}
}
