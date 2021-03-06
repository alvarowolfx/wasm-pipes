package sink

import (
	"bytes"
	"context"
	"net/http"

	"github.com/apex/log"
	"gocloud.dev/pubsub"
)

type httpSink struct {
	uri   string
	input *pubsub.Subscription
}

func NewHTTPSink(deps SinkDeps) (Sink, error) {
	return &httpSink{uri: deps.URI, input: deps.Input}, nil
}

func (s *httpSink) Start() {
	go s.listenMessages()
}

func (s *httpSink) listenMessages() {
	if s.input == nil {
		return
	}

	for {
		ctx := context.Background()
		msg, err := s.input.Receive(ctx)
		if err != nil {
			log.Errorf("relaying message to sink: %v", err)
			break
		}
		log.Infof("got message to send to sink: %q", msg.Body)

		go func(msg *pubsub.Message) {
			client := &http.Client{}
			req, err := http.NewRequest("POST", s.uri, bytes.NewReader(msg.Body))
			if err != nil {
				log.Errorf("relaying message to sink: %v", err)
				msg.Ack()
				return
			}
			res, err := client.Do(req)
			if err != nil {
				log.Errorf("relaying message to sink: %v", err)
				msg.Ack()
				return
			}
			log.Infof("response from sink: %s", res.Status)
			msg.Ack()
		}(msg)
	}
}

func (s *httpSink) Shutdown(ctx context.Context) error {
	return nil
}
