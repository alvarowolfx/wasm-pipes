package source

import (
	"context"
	"strings"

	"gocloud.dev/pubsub"
)

type Source interface {
	Start()
	Stop(context.Context) error
}

type SourceDeps struct {
	URI           string
	EventBusTopic *pubsub.Topic
}

func NewSource(deps SourceDeps) (Source, error) {
	if strings.HasPrefix(deps.URI, "http") {
		return NewHTTPSource(deps)
	}
	return NewPubsubSource(deps)
}
