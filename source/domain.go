package source

import (
	"context"
	"strings"

	"gocloud.dev/pubsub"
)

type Source interface {
	Start()
	Shutdown(context.Context) error
}

type SourceDeps struct {
	URI    string
	Output *pubsub.Topic
}

func NewSource(deps SourceDeps) (Source, error) {
	if strings.HasPrefix(deps.URI, "http") {
		return NewHTTPSource(deps)
	}
	return NewPubsubSource(deps)
}
