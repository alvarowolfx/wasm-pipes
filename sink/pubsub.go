package sink

import (
	"context"

	"github.com/apex/log"
	"gocloud.dev/pubsub"
)

type pubsubSink struct {
	uri    string
	input  *pubsub.Subscription
	output *pubsub.Topic
}

func NewPubsubSink(deps SinkDeps) (Sink, error) {
	return &pubsubSink{uri: deps.URI, input: deps.Input}, nil
}

func (s *pubsubSink) Start() {
	ctx := context.Background()
	var err error
	s.output, err = pubsub.OpenTopic(ctx, s.uri)
	if err != nil {
		log.Fatalf("err opening pubsub sink topic: %v", err)
	}

	go s.listenMessages()
}

func (s *pubsubSink) listenMessages() {
	if s.input == nil {
		return
	}

	for {
		ctx := context.Background()
		msg, err := s.input.Receive(ctx)
		if err != nil {
			log.Errorf("receiving message on sink: %v", err)
			break
		}
		log.Infof("got message on sink: %q", msg.Body)

		go func(msg *pubsub.Message) {
			serr := s.output.Send(ctx, &pubsub.Message{
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

func (s *pubsubSink) Shutdown(ctx context.Context) error {
	if s.output != nil {
		return s.output.Shutdown(ctx)
	}
	return nil
}
