# AWS SQS Sdk & Golang Complete Cheat Sheet


# SQS (Simple Queue Service) Query Examples

This page should help you to understand how to use AWS SQS with golang using the official [aws-sdk-go-v2/service/sqs](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/sqs).

We will start with some basics like how to create the client to rather more complex operations such as Sending Message etc.

## Setup

Connecting Go application to SQS Client is quite easy. You just need to load your config and create the client.

You need to get the below packages
1. go get -u `github.com/aws/aws-sdk-go-v2/config`
1. go get -u `github.com/aws/aws-sdk-go-v2/service/sqs`

{{< highlight go "linenos=false" >}}
package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
}
{{< / highlight >}}

## Create Queue

AS we know, the messages are stored in the queues, so without a queue we will not be able to operate any function. Creating queues from AWS CLI is quite easy, it is also easy with AWS SDK. Let's create a queue with the name `test_queue`.

{{< highlight go "linenos=false" >}}
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queue, err := sqsClient.CreateQueue(context.TODO(), &sqs.CreateQueueInput{
		QueueName:  aws.String("test_queue"),
		Attributes: nil,
		Tags:       nil,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("the queue url is %v", *queue.QueueUrl)
}
{{< / highlight >}}

Now we have a queue with the name `test_queue`.

## Fetch Queue URL

Some operations such as sending message to a queue or deleting the queue needs the queue url as input. Even though we know the name of the queue, we still need to fetch the URL of it. It is quite easy.

{{< highlight go "linenos=false" >}}
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queue, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName:              aws.String("test_queue"),
		QueueOwnerAWSAccountId: nil,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("the queue url is %v", *queue.QueueUrl)
}
{{< / highlight >}}

## Delete Queue

Deleting the queue is also quite easy. You just need to call the `DeleteQueue` function. You will be using this a lot if you create temporary queues and need to delete them after a while.

{{< highlight go "linenos=false" >}}
package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queue, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName:              aws.String("test_queue"),
		QueueOwnerAWSAccountId: nil,
	})
	if err != nil {
		panic(err)
	}

	_, err = sqsClient.DeleteQueue(context.TODO(), &sqs.DeleteQueueInput{QueueUrl: queue.QueueUrl})
	if err != nil {
		panic(err)
	}
}
{{< / highlight >}}

## List Queues

Sometimes we need to list the queues to see which queues there are in the AWS. You can even set a prefix to get the queues with wanted prefix.

Let's fetch the queues with prefix `test`.

{{< highlight go "linenos=false" >}}
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queues, err := sqsClient.ListQueues(context.TODO(), &sqs.ListQueuesInput{
		MaxResults:      nil,
		NextToken:       nil,
		QueueNamePrefix: aws.String("test"),
	})
	if err != nil {
		panic(err)
	}

	for _, queueUrl := range queues.QueueUrls {
		log.Println(queueUrl)
	}
}
{{< / highlight >}}

## Send Message

The most important part of the queues is mostly is receiving & sending messages part. The whole point of the queue is to become the middle man which has the responsibility of taking and delivering the messages. Sending a message in SQS quite easy.


{{< highlight go "linenos=false" >}}
package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type YourStruct struct {
	University string
	Major      string
	Level      string
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	data := YourStruct{
		University: "Technical University of Munich",
		Major:      "Informatics",
		Level:      "Graduate",
	}

	bytes, _ := json.Marshal(&data)

	queue, _ := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String("test_queue"),
	})

	res, err := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody:             aws.String(string(bytes)),
		QueueUrl:                queue.QueueUrl,
		DelaySeconds:            0,
		MessageAttributes:       nil,
		MessageDeduplicationId:  nil,
		MessageGroupId:          nil,
		MessageSystemAttributes: nil,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("the message with id %v is sent", *res.MessageId)
}
{{< / highlight >}}

## Receive Message

The other most important part of the queues are receiving the message, because as we know the messages are only sent for someone to read it. I think that could be a good quote so let's make a one.

{{< admonition quote "Hasan Ocak" >}}
Messages are only sent for someone to read it.
{{< /admonition >}}

Handling messages in SQS can be quite complex. If you want to learn how to receive message and handle them in parallel, check out my other post [Golang Sqs Consumer Worker Pool](https://ocakhasan.github.io/golang-sqs-consumer-worker-pool/).

{{< highlight go "linenos=false" >}}
package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queue, _ := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String("test_queue"),
	})

	messages, err := sqsClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:                queue.QueueUrl,
		AttributeNames:          nil,
		MaxNumberOfMessages:     10, // max is 10
		MessageAttributeNames:   nil,
		ReceiveRequestAttemptId: nil,
		VisibilityTimeout:       0,
		WaitTimeSeconds:         0,
	})
	if err != nil {
		panic(err)
	}

	for _, message := range messages.Messages {
		log.Printf("the message body is %v", *message.Body)
	}
}
{{< / highlight >}}

Output will be something like this

```
2023/10/18 22:22:57 the message body is {"University":"Technical University of Munich","Major":"Informatics","Level":"Graduate"}
```

### How to Decode the Message Into Your Struct

You just need to decode it via `json.Unmarshal` to your struct.

```go
for _, message := range messages.Messages {
    var data YourStruct
    _ = json.Unmarshal([]byte(*message.Body), &data)

    log.Printf("the received data is %+v", data)
}
```

## Delete Message

In general after handling the message, you will want to delete the message so it will not be received and handled again. Of course this can change according to your needs. Deleting the message is quite easy.

{{< highlight go "linenos=false" >}}
package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type YourStruct struct {
	University string
	Major      string
	Level      string
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queue, _ := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String("test_queue"),
	})

	// It can be read from message variable after receiving from SQS
	messageHandle := "YOUR_MESSAGE_HANDLE"

	_, err = sqsClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      queue.QueueUrl,
		ReceiptHandle: &messageHandle,
	})
	if err != nil {
		panic(err)
	}
}
{{< / highlight >}}

## Purge Queue

If you want to delete all of the messages in the queue, you need to use the `PurgeQueue` method. 

{{< highlight go "linenos=false" >}}
package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	sqsClient := sqs.NewFromConfig(cfg)

	queue, _ := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String("test_queue"),
	})

	_, err = sqsClient.PurgeQueue(context.TODO(), &sqs.PurgeQueueInput{
		QueueUrl: queue.QueueUrl,
	})
	if err != nil {
		panic(err)
	}
}
{{< / highlight >}}

## Motivation

This post is too basic with SQS but it might help a beginner programmer who is trying to do some stuff with SQS. While I was doing some development with `DynamoDB`, I got help from this [blog](https://dynobase.dev/dynamodb-golang-query-examples/). It really helped me to get the development going.

I wanted to do it the same with SQS. The reason I do it with SQS is I work with it often, so I know some stuff about it.

Hope, it will help you. If it helps or not, you can reach out to me.

Let's end the topic with a calming music. Enjoy :musical_note: :computer:

{{< music "https://music.163.com/#/song?id=1322030323" >}}

