package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gocloud.dev/pubsub"
	"gocloud.dev/server"
)

type httpSource struct {
	path   string
	port   string
	output *pubsub.Topic
	s      *server.Server
}

func NewHTTPSource(deps SourceDeps) (Source, error) {
	uri := deps.URI
	url, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("err invalid source uri: %w", err)
	}
	port := url.Port()
	if port == "" {
		port = "80"
	}
	path := url.Path
	return &httpSource{port: port, path: path, output: deps.Output}, nil
}

func (r *httpSource) Start() {
	r.s = server.New(http.DefaultServeMux, nil)

	http.HandleFunc(r.path, func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		metadata := map[string]string{}
		for h, v := range req.Header {
			metadata[h] = strings.Join(v, ",")
		}
		if req.Method == http.MethodPost || req.Method == http.MethodPut {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				res.WriteHeader(400)
				fmt.Fprintln(res, err.Error())
				return
			}
			err = r.output.Send(ctx, &pubsub.Message{
				Body:     body,
				Metadata: metadata,
			})
			if err != nil {
				res.WriteHeader(500)
				fmt.Fprintln(res, err.Error())
				return
			}
			res.WriteHeader(200)
			fmt.Fprintln(res, "ok")
		} else if req.Method == http.MethodGet {
			body, err := json.Marshal(req.URL.Query())
			if err != nil {
				res.WriteHeader(400)
				fmt.Fprintln(res, err.Error())
				return
			}
			err = r.output.Send(ctx, &pubsub.Message{
				Body:     body,
				Metadata: metadata,
			})
			if err != nil {
				res.WriteHeader(500)
				fmt.Fprintln(res, err.Error())
				return
			}
			res.WriteHeader(200)
			fmt.Fprintln(res, "ok")
		} else {
			res.WriteHeader(400)
			fmt.Fprintln(res, "invalid method")
		}
	})

	log.Printf("Starting source on port %s", r.port)
	if err := r.s.ListenAndServe(fmt.Sprintf(":%s", r.port)); err != nil {
		log.Fatalf("%v", err)
	}
}

func (r *httpSource) Shutdown(ctx context.Context) error {
	if r.s != nil {
		return r.s.Shutdown(ctx)
	}
	return nil
}
