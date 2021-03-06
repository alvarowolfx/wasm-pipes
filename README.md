# WASM Pipes

<img src="./.images/logo.png" width="64" height="64" />

🚧**THIS IS A WORK IN PROGRESS - DON'T TRY IN PRODUCTION**🚧

A lot of work on stream processing is doing some "data ping-pong", which is to relay data from one source to another, doing some small transformations or filtering.

> ~60% of stream processing is spent doing mundane transformation tasks like format unification for ML workloads, filtering for privacy, simple enrichments like geo-ip translations, etc. Yet the baseline complexity to stand up something “simple” often involves 3 or 4 distributed systems, multiple nights reading configurations and man-pages, and a few too many shots of espresso to start seeing the future. And once you are done, you end up ping-ponging the data back and forth between storage and compute, when all you had to do was to remove a field out of a JSON object. To the data engineer, it feels like an endless game of system whack-a-mole just to start to do the interesting work of actually understanding the data.
> From [Vectorized article](https://vectorized.io/blog/wasm-architecture/)

The idea behind this project is to allows developers to create easily customizable Message Queue pipelines using Web Assembly to transform/filter data. On top of that, also given access to some services, like Logs, Cache and routing to message queues.

![Architecture](./.images/architecture.png)

[Demo Video](https://www.loom.com/share/e60d4850342145818b3606a324c38643)

## Build/Install

You need go installed to build this project.

```
> go install ./ # or make install
> wasm-pipes
```

## Samples

- [Relaying data from a Webhook to a NATS Broker](./demos/relay-webhook-to-nats.md)
- [Transform Pipeline - Uppercase](./demos/transform-uppercase-pipeline.md)
- [Dedup Sequence Pipeline](./demos/dedup-sequence-pipeline.md)
- [Filtering/Routing Pipeline](./demos/filtering-pipeline.md)

## Building Web Assembly Modules

🚧**Work in Progress**🚧

## Supported Sources and Sinks

Every pipeline item have a entry point (Source) and a output point (Sink). Here is the list of supported sources and sinks that can be used on the project. Most of the supported sources/sinks here are from the [Go Cloud](https://gocloud.dev) project, currently the only addition are the HTTP ones.

### Source

Configurable by the env variable `SOURCE_URI`.

- HTTP
  - Example URI - `http://:9000/events`
  - Open an HTTP endpoint on the given port and path to receive data.
    - POST/PUT requests will relay the request body through the pipeline.
    - GET requests will relay Query Params encoded in JSON through the pipeline.
- GCP PubSub - [gocloud docs](https://gocloud.dev/howto/pubsub/subscribe/#gcp)
  - Example URI - `gcppubsub://projects/my-project/subscriptions/my-subscription`
  - Extra required env variables - from Go Cloud Docs
    - Will use Application Default Credentials; if you have authenticated via `gcloud auth application-default login`, it will use those credentials. See [Application Default Credentials](https://cloud.google.com/docs/authentication/production) to learn about authentication alternatives, including using environment variables.
- AWS SQS - [gocloud docs](https://gocloud.dev/howto/pubsub/subscribe/#sqs)
  - Example URI - `awssqs://sqs.us-east-2.amazonaws.com/123456789012/"+ "myqueue?region=us-east-2`
  - Extra required env variables
    - Will create a default AWS Session with the SharedConfigEnable option enabled; if you have authenticated with the AWS CLI, it will use those credentials. See [AWS Session](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/) to learn about authentication alternatives, including using environment variables.
- Azure Service Bus - [gocloud docs](https://gocloud.dev/howto/pubsub/subscribe/#azure)
  - Example URI - `azuresb://mytopic?subscription=mysubscription`
  - Extra required env variables
    - Will use the environment variable `SERVICEBUS_CONNECTION_STRING` to obtain the Service Bus Connection String you need to copy [from the Azure portal](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-how-to-use-topics-subscriptions#get-the-connection-string).
- RabbitMQ - [gocloud docs](https://gocloud.dev/howto/pubsub/subscribe/#rabbitmq)
  - Example URI - `rabbit://myqueue`
  - Extra required env variables
    - The RabbitMQ’s server is discovered from the `RABBIT_SERVER_URL` environment variable (which is something like `amqp://guest:guest@localhost:5672/`).
- NATS - [gocloud docs](https://gocloud.dev/howto/pubsub/subscribe/#nats)
  - Example URI - `nats://example.mysubject?queue=`
  - Extra required env variables
    - The NATS server is discovered from the `NATS_SERVER_URL` environment variable (which is something like `nats://nats.example.com`).
- Kafka - [gocloud docs](https://gocloud.dev/howto/pubsub/subscribe/#kafka)
  - Example URI - `kafka://my-group?topic=my-topic`
  - Extra required env variables
    - The brokers in the Kafka cluster are discovered from the `KAFKA_BROKERS` environment variable (which is a comma-delimited list of hosts, something like `1.2.3.4:9092,5.6.7.8:9092`).

### Sink

Configurable by the env variable `SINK_URI`.

- HTTP
  - Example URI - `http://example.com:9000/webhook`
  - Will `POST` the pipeline output to the configure URL as an HTTP request.
- GCP PubSub - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#gcp)
  - Example URI - `gcppubsub://projects/myproject/topics/mytopic`
  - Extra required env variables - from Go Cloud Docs
    - Will use Application Default Credentials; if you have authenticated via `gcloud auth application-default login`, it will use those credentials. See [Application Default Credentials](https://cloud.google.com/docs/authentication/production) to learn about authentication alternatives, including using environment variables.
- AWS SNS - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#sns)
  - Example URI - `awssns:///arn:aws:sns:us-east-2:123456789012:mytopic?region=us-east-2"`
  - Extra required env variables
    - Will create a default AWS Session with the SharedConfigEnable option enabled; if you have authenticated with the AWS CLI, it will use those credentials. See [AWS Session](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/) to learn about authentication alternatives, including using environment variables.
- AWS SQS - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#sqs)
  - Example URI - `myqueueawssqs://sqs.us-east-2.amazonaws.com/123456789012/?region=us-east-2`
  - Extra required env variables
    - Will create a default AWS Session with the SharedConfigEnable option enabled; if you have authenticated with the AWS CLI, it will use those credentials. See [AWS Session](https://docs.aws.amazon.com/sdk-for-go/api/aws/session/) to learn about authentication alternatives, including using environment variables.
- Azure Service Bus - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#azure)
  - Example URI - `azuresb://mytopic`
  - Extra required env variables
    - Will use the environment variable `SERVICEBUS_CONNECTION_STRING` to obtain the Service Bus Connection String you need to copy [from the Azure portal](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-dotnet-how-to-use-topics-subscriptions#get-the-connection-string).
- RabbitMQ - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#rabbitmq)
  - Example URI - `rabbit://myexchange`
  - Extra required env variables
    - The RabbitMQ’s server is discovered from the `RABBIT_SERVER_URL` environment variable (which is something like `amqp://guest:guest@localhost:5672/`).
- NATS - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#nats)
  - Example URI - `nats://example.mysubject`
  - Extra required env variables
    - The NATS server is discovered from the `NATS_SERVER_URL` environment variable (which is something like `nats://nats.example.com`).
- Kafka - [gocloud docs](https://gocloud.dev/howto/pubsub/publish/#kafka)
  - Example URI - `kafka://my-topic`
  - Extra required env variables
    - The brokers in the Kafka cluster are discovered from the `KAFKA_BROKERS` environment variable (which is a comma-delimited list of hosts, something like `1.2.3.4:9092,5.6.7.8:9092`).

## References

- About Data Ping-Pong - https://vectorized.io/blog/wasm-architecture/
- Go Cloud Development Kit - https://gocloud.dev/
- WASM Engine - https://github.com/wasmerio/wasmer-go
