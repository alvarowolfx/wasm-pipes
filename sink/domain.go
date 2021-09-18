package sink

import (
	"context"
	"strings"

	"gocloud.dev/pubsub"
)

type Sink interface {
	Start()
	Stop(context.Context) error
}

type SinkDeps struct {
	URI         string
	EventBusSub *pubsub.Subscription
}

func NewSink(deps SinkDeps) (Sink, error) {
	if strings.HasPrefix(deps.URI, "http") {
		return NewHTTPSink(deps)
	}
	return NewPubsubSink(deps)
}
