package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvarowolfx/wasm-pipes/config"
	"github.com/alvarowolfx/wasm-pipes/engine"
	"github.com/alvarowolfx/wasm-pipes/infra"
	"github.com/alvarowolfx/wasm-pipes/sink"
	"github.com/alvarowolfx/wasm-pipes/source"
	"github.com/apex/log"
	"github.com/joho/godotenv"
	"gocloud.dev/blob"
	"gocloud.dev/pubsub"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warn("Error loading .env file")
	}
	cfg := config.LoadConfigFromEnv()
	fmt.Println("config:", cfg)
	infra.LoadProviders()

	ctx := context.Background()

	cacheStore, err := infra.NewCache("mem://")
	checkErr("could not create engine cache", err)

	wasmBucket, err := blob.OpenBucket(ctx, cfg.WasmBucketURI)
	checkErr("could not open engine bucket", err)
	defer func() {
		cerr := wasmBucket.Close()
		if cerr != nil {
			log.Errorf("could not close engine bucket: %v", cerr)
		}
	}()

	sourceURI := "mem://sourceTopic"
	sourceOutputTopic, err := pubsub.OpenTopic(ctx, sourceURI)
	checkErr("could not create internal source topic", err)
	defer handleStoppable(sourceOutputTopic)
	sourceOutputSub, err := pubsub.OpenSubscription(ctx, sourceURI)
	checkErr("could not create internal source subscription", err)
	defer handleStoppable(sourceOutputSub)

	engineURI := "mem://engineTopic"
	engineOutputTopic, err := pubsub.OpenTopic(ctx, engineURI)
	checkErr("could not create internal engine topic", err)
	defer handleStoppable(engineOutputTopic)
	engineOutputSub, err := pubsub.OpenSubscription(ctx, engineURI)
	checkErr("could not create internal engine subscription", err)
	defer handleStoppable(engineOutputSub)

	src, err := source.NewSource(source.SourceDeps{
		URI:    cfg.SourceURI,
		Output: sourceOutputTopic,
	})
	checkErr("could not create source", err)
	go src.Start()
	defer handleStoppable(src)

	eng, err := engine.NewEngine(&engine.EngineDeps{
		Filename: cfg.WasmFilename,
		Bucket:   wasmBucket,
		Cache:    cacheStore,
		Input:    sourceOutputSub,
		Output:   engineOutputTopic,
	})
	checkErr("could not create wasm engine", err)
	go eng.Start()
	defer handleStoppable(eng)

	sk, err := sink.NewSink(sink.SinkDeps{
		URI:   cfg.SinkURI,
		Input: engineOutputSub,
	})
	checkErr("could not create sink", err)
	go sk.Start()
	defer handleStoppable(sk)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("wasm-pipes started")
	<-done
	log.Info("wasm-pipes received signal to stop")
	log.Info("wasm-pipes stopped")
}

func checkErr(msg string, err error) {
	if err != nil {
		log.Errorf("%s: %v", msg, err)
		os.Exit(1)
	}
}

type stoppable interface {
	Shutdown(context.Context) error
}

func handleStoppable(s stoppable) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	serr := s.Shutdown(ctx)
	if serr != nil {
		log.Errorf("could not stop source: %v", serr)
	}
}
