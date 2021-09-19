package engine

import (
	"context"

	"github.com/apex/log"
	"github.com/wasmerio/wasmer-go/wasmer"
	"gocloud.dev/blob"
	"gocloud.dev/pubsub"
)

type engine struct {
	filename string
	input    *pubsub.Subscription
	output   *pubsub.Topic
	bucket   *blob.Bucket
	module   *wasmer.Module
	store    *wasmer.Store
	deps     *EngineDeps
}

func NewEngine(deps *EngineDeps) (Engine, error) {
	e := &engine{
		deps:     deps,
		filename: deps.Filename,
		bucket:   deps.Bucket,
		input:    deps.Input,
		output:   deps.Output,
	}
	err := e.init()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *engine) init() error {
	wasmEng := wasmer.NewEngine()
	e.store = wasmer.NewStore(wasmEng)
	filename := "main.wasm"
	if e.filename != "" {
		filename = e.filename
	}
	wasmBytes, err := e.bucket.ReadAll(context.Background(), filename)
	/*wasmF, err := os.Open(e.uri)
	if err != nil {
		log.Errorf("failed to open wasm: %v", err)
		return err
	}

	wasmBytes, err := ioutil.ReadAll(wasmF)*/
	if err != nil {
		log.Errorf("failed to read wasm: %v", err)
		return err
	}

	module, err := wasmer.NewModule(e.store, wasmBytes)
	if err != nil {
		log.Fatalf("failed to compile module: %v", err)
		return err
	}

	e.module = module
	return nil
}

func (e *engine) Start() {
	go e.listenMessages()
}

func (e *engine) listenMessages() {
	if e.input == nil {
		return
	}

	for {
		ctx := context.Background()
		msg, err := e.input.Receive(ctx)
		if err != nil {
			log.Errorf("receiving message on engine: %v", err)
			break
		}
		log.Infof("got message on engine: %q", msg.Body)

		go func(msg *pubsub.Message) {
			rerr := e.runWithMsg(msg)
			if rerr != nil {
				log.Errorf("failed to call the `run` function: %v", err)
				msg.Ack()
				return
			}

			log.Info("event processed successfully")
			msg.Ack()
		}(msg)
	}
}

func (e *engine) runWithMsg(msg *pubsub.Message) error {
	var err error
	var instance *wasmer.Instance
	importObject := wasmer.NewImportObject()

	execVars := &ExecVars{
		Message: msg,
	}

	log.Debugf("found %d wasm plugins", len(plugins))
	for pName, p := range plugins {
		err = p.Init(e.deps, e.store, execVars)
		if err != nil {
			log.Errorf("failed to init plugin %s: %v", pName, err)
			return err
		}
		objects := p.ImportedObjects()
		log.Debugf("registering plugin %s", pName)
		importObject.Register(pName, objects)
	}

	instance, err = wasmer.NewInstance(e.module, importObject)
	execVars.Instance = instance
	if err != nil {
		log.Errorf("failed to instantiate the module: %v", err)
		return err
	}

	run, err := instance.Exports.GetFunction("run")
	if err != nil {
		log.Errorf("failed to get the `run` function: %v", err)
		return err
	}

	_, err = run()
	return err
}

func (e *engine) Shutdown(ctx context.Context) error {
	return nil
}
