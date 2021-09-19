package source

import (
	"context"

	"github.com/apex/log"

	"gocloud.dev/pubsub"
)

type pubsubSource struct {
	uri    string
	output *pubsub.Topic
	input  *pubsub.Subscription
}

func NewPubsubSource(deps SourceDeps) (Source, error) {
	return &pubsubSource{uri: deps.URI, output: deps.Output}, nil
}

func (r *pubsubSource) Start() {
	ctx := context.Background()
	var err error
	r.input, err = pubsub.OpenSubscription(ctx, r.uri)
	if err != nil {
		log.Fatalf("err opening pubsub subscription: %v", err)
	}

	go r.listenMessages()
}

func (r *pubsubSource) listenMessages() {
	if r.input == nil {
		return
	}

	for {
		ctx := context.Background()
		msg, err := r.input.Receive(ctx)
		if err != nil {
			log.Errorf("receiving message: %v", err)
			break
		}
		log.Infof("got message: %q\n", msg.Body)

		go func(msg *pubsub.Message) {
			serr := r.output.Send(ctx, &pubsub.Message{
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

func (r *pubsubSource) Shutdown(ctx context.Context) error {
	if r.input != nil {
		return r.input.Shutdown(ctx)
	}
	return nil
}
