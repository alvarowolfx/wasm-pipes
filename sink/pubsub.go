package sink

import (
	"context"

	"github.com/apex/log"
	"gocloud.dev/pubsub"
)

type pubsubSink struct {
	uri      string
	eventBus *pubsub.Subscription
	topic    *pubsub.Topic
}

func NewPubsubSink(deps SinkDeps) (Sink, error) {
	return &pubsubSink{uri: deps.URI, eventBus: deps.EventBusSub}, nil
}

func (s *pubsubSink) Start() {
	ctx := context.Background()
	var err error
	s.topic, err = pubsub.OpenTopic(ctx, s.uri)
	if err != nil {
		log.Fatalf("err opening pubsub sink topic: %v", err)
	}

	go s.listenMessages()
}

func (s *pubsubSink) listenMessages() {
	if s.topic == nil {
		return
	}

	for {
		ctx := context.Background()
		msg, err := s.eventBus.Receive(ctx)
		if err != nil {
			log.Errorf("receiving message on sink: %v", err)
			break
		}
		log.Infof("got message on sink: %q", msg.Body)

		go func(msg *pubsub.Message) {
			serr := s.topic.Send(ctx, &pubsub.Message{
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

func (s *pubsubSink) Stop(ctx context.Context) error {
	if s.topic != nil {
		return s.topic.Shutdown(ctx)
	}
	return nil
}
