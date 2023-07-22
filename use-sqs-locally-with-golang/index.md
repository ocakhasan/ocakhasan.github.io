# Local SQS Setup With Golang


# Introduction

Before going any further all of the code can be found in the [local-go-sqs-setup](https://github.com/ocakhasan/local-sqs-go-setup).

Welcome to our blog post on using local Amazon Simple Queue Service (SQS) with Golang! As developers, we understand the importance of reliable and efficient message queuing systems for building scalable and decoupled applications. Amazon SQS, a fully managed message queuing service, is a popular choice among cloud-based solutions. However, during the development and testing phase, relying on the cloud-based SQS might not always be the best option due to network latency and potential costs.

In this comprehensive guide, we will dive into the world of local SQS and explore how to harness its power within your Golang applications. By setting up and running SQS on your local environment, you can achieve a seamless development and testing experience, all while reducing the dependency on the cloud infrastructure.

We will walk you through the step-by-step process of configuring your local SQS instance, integrating it with your Golang application, and simulating real-world scenarios. Whether you're building a microservices architecture or need a reliable message queue for inter-process communication, local SQS can be a game-changer in your development workflow.

## Technologies we use

1. Docker
2. We will use the docker image of the [softwaremill/elasticmq](https://github.com/softwaremill/elasticmq) for the local SQS
3. Golang

`softwaremill/elasticmq` is a tool which runs a SQS compatible server on your local environment.


## Setup the Local SQS

To be able to run an sqs server locally, simply run the command

```bash
docker run -p 9324:9324 -p 9325:9325 softwaremill/elasticmq
```

Then, go to your browser and enter the url of `http://localhost:9325`. You will see something like this.

![Browser image](../../images/sqs.png)

## Create A Queue

To create a queue, run the command

```bash
aws sqs create-queue --endpoint-url http://localhost:9324 --queue-name test_queue --region eu-west-1
```

the response will be something like this.

```json
{
    "QueueUrl": "http://localhost:9324/000000000000/test_queue"
}
```

And the browser will be 

![SQS Create Local Queue](../../images/sqs_create_queue.png)

## How To Integrate With Go

Normally that's how you create a sqs client in go to list the queue urls.

```go
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	res, err := sqsClient.ListQueues(context.Background(), &sqs.ListQueuesInput{
		MaxResults:      aws.Int32(10),
		NextToken:       nil,
		QueueNamePrefix: nil,
	})

	if err != nil {
		log.Printf("error while listing the queues")
	}

	for _, queue := range res.QueueUrls {
		log.Println(queue)
	}
}
```

However, to be able to connect to local queue we need an [EndpointResolverWithOptions](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/) which will redirect the requests to `http://localhost:9324`

```go
type EndpointResolverWithOptions interface {
	ResolveEndpoint(service, region string, options ...interface{}) (Endpoint, error)
}
```

To be able to do it, we can create a simple struct which implements the `EndpointResolverWithOptions` interface. 

```go
type localResolver struct{}

func (l localResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL:           "http://localhost:9324",
		SigningRegion: "eu-west-1",
	}, nil
}
```

The all of the code is

```go
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type localResolver struct{}

func (l localResolver) ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL:           "http://localhost:9324",
		SigningRegion: "eu-west-1",
	}, nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	cfg.EndpointResolverWithOptions = localResolver{}

	sqsClient := sqs.NewFromConfig(cfg)

	res, err := sqsClient.ListQueues(context.Background(), &sqs.ListQueuesInput{
		MaxResults:      aws.Int32(10),
		NextToken:       nil,
		QueueNamePrefix: nil,
	})

	if err != nil {
		log.Printf("error while listing the queues")
	}

	for _, queue := range res.QueueUrls {
		log.Println(queue)
	}
}
```

When you run this code with the command 

```bash
go run main.go
```

You will see something like

```
2023/07/22 20:27:20 http://localhost:9324/000000000000/test_queue
```

## REFERENCES

- https://github.com/softwaremill/elasticmq
- https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs

Thanks for reading, please let me know if you have any questions.
