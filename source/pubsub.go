package source

import (
	"context"

	"github.com/apex/log"

	"gocloud.dev/pubsub"
)

type pubsubSource struct {
	uri      string
	eventBus *pubsub.Topic
	sub      *pubsub.Subscription
}

func NewPubsubSource(deps SourceDeps) (Source, error) {
	return &pubsubSource{uri: deps.URI, eventBus: deps.EventBusTopic}, nil
}

func (r *pubsubSource) Start() {
	ctx := context.Background()
	var err error
	r.sub, err = pubsub.OpenSubscription(ctx, r.uri)
	if err != nil {
		log.Fatalf("err opening pubsub subscription: %v", err)
	}

	go r.listenMessages()
}

func (r *pubsubSource) listenMessages() {
	if r.sub == nil {
		return
	}

	for {
		ctx := context.Background()
		msg, err := r.sub.Receive(ctx)
		if err != nil {
			log.Errorf("receiving message: %v", err)
			break
		}
		log.Infof("got message: %q\n", msg.Body)

		go func(msg *pubsub.Message) {
			serr := r.eventBus.Send(ctx, &pubsub.Message{
				Body:     msg.Body,
				Metadata: msg.Metadata,
			})
			if serr != nil {
				log.Errorf("sending message to msg bus: %v", serr)
				msg.Nack()
				return
			}
			msg.Ack()
		}(msg)
	}
}

func (r *pubsubSource) Stop(ctx context.Context) error {
	if r.sub != nil {
		return r.sub.Shutdown(ctx)
	}
	return nil
}
