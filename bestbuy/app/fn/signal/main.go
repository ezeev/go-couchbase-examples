package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/ezeev/go-couchbase-examples/bestbuy/app/model"
	"github.com/golang/protobuf/proto"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, name events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer
	//serialize
	var signal model.Signal
	err := json.Unmarshal([]byte(name.Body), &signal)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	// if timestamp is null, populate it with current time
	if signal.ClickTime == 0 {
		signal.ClickTime = time.Now().Unix()
	}

	fmt.Printf("deserialized struct: %v", signal)

	// marshall to protobuff
	message, err := proto.Marshal(&signal)
	if err != nil {
		fmt.Printf("Error mashalling to proto: %v", err)
	}

	// put to kinesis stream
	fmt.Print(message)
	s := session.New(&aws.Config{Region: aws.String(os.Getenv("AZ"))})
	kc := kinesis.New(s)
	streamName := aws.String(os.Getenv("SIGNAL_STREAM"))
	entry := kinesis.PutRecordInput{}
	entry.StreamName = streamName
	entry.SetData(message)
	entry.PartitionKey = aws.String(fmt.Sprintf("%d", time.Now().Second()))

	out, err := kc.PutRecord(&entry)
	if err != nil {
		fmt.Printf("error putting kinesis record: %v", err)
	}
	fmt.Printf("Recived kinesis response to put: %v", out)

	body, err := json.Marshal(map[string]interface{}{
		"message": "Received signal",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
