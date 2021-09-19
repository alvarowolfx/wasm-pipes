package infra

import (
	"github.com/apex/log"

	/* Load all pubsub providers */
	"gocloud.dev/blob"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/awssnssqs"
	_ "gocloud.dev/pubsub/azuresb"
	_ "gocloud.dev/pubsub/gcppubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
	_ "gocloud.dev/pubsub/mempubsub"
	_ "gocloud.dev/pubsub/natspubsub"
	_ "gocloud.dev/pubsub/rabbitpubsub"

	/* Load all blob providers */
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"

	_ "github.com/alvarowolfx/wasm-pipes/engine/services"
)

func LoadProviders() {
	log.Infof("loading providers")
	log.Infof("blob storage providers: %v", blob.DefaultURLMux().BucketSchemes())
	log.Infof("pubsub topic providers: %v", pubsub.DefaultURLMux().TopicSchemes())
	log.Infof("pubsub subscription providers: %v", pubsub.DefaultURLMux().SubscriptionSchemes())
	log.Info("cache providers: [mem]")
}
