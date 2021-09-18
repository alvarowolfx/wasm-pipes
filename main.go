package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alvarowolfx/wasm-pipes/config"
	"github.com/alvarowolfx/wasm-pipes/infra"
	"github.com/alvarowolfx/wasm-pipes/sink"
	"github.com/alvarowolfx/wasm-pipes/source"
	"github.com/apex/log"
	"github.com/joho/godotenv"
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
	eventBusURI := "mem://eventbus"
	eventBusTopic, err := pubsub.OpenTopic(ctx, eventBusURI)
	if err != nil {
		log.Fatalf("could create internal event bus topic: %v", err)
		return
	}
	eventBusSub, err := pubsub.OpenSubscription(ctx, eventBusURI)
	if err != nil {
		log.Fatalf("could create internal event bus sub: %v", err)
		return
	}

	src, err := source.NewSource(source.SourceDeps{
		URI:           cfg.SourceURI,
		EventBusTopic: eventBusTopic,
	})
	if err != nil {
		log.Fatalf("could not create source: %v", err)
		return
	}
	go src.Start()
	defer handleStoppable(src)

	sk, err := sink.NewSink(sink.SinkDeps{
		URI:         cfg.SinkURI,
		EventBusSub: eventBusSub,
	})
	if err != nil {
		log.Fatalf("could not create sink: %v", err)
		return
	}
	go sk.Start()
	defer handleStoppable(sk)

	/*go func() {
		for {
			ctx := context.Background()
			msg, err := eventBusSub.Receive(ctx)
			if err != nil {
				log.Errorf("receiving message on event bus: %v", err)
				break
			}
			log.Infof("got message on event bus: %q", msg.Body)
			msg.Ack()
		}
	}()*/

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("wasm-pipes started")
	<-done
	log.Info("wasm-pipes received signal to stop")
	log.Info("wasm-pipes stopped")
}

type stoppable interface {
	Stop(context.Context) error
}

func handleStoppable(s stoppable) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	serr := s.Stop(ctx)
	if serr != nil {
		log.Errorf("could not stop source: %v", serr)
	}
}
