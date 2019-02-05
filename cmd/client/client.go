package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/afiskon/golang-jaeger-example"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/trace"
)

func main() {
	utils.InitJaeger("client")
	for {
		ctx, span := trace.StartSpan(context.Background(), "main")
		span.AddAttributes(
			trace.StringAttribute("method", "GET"),
		)
		log.Printf("TraceId: %s\n", span.SpanContext().TraceID.String())
		resp, err := sendPostRequest(ctx, "http://localhost:8080/api/v1/ping")
		if err != nil {
			log.Printf("sendPostRequest: %v\n", err)
		}
		fmt.Printf("Response: %s\n", resp)
		span.End()

		time.Sleep(5 * time.Second)
	}
}

func sendPostRequest(ctx context.Context, url string) (string, error) {
	_ /* ctx2 */, span := trace.StartSpan(ctx, "sendPostRequest")
	defer span.End()

	spanContextJson, err := json.Marshal(span.SpanContext())
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url,  nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Span-Context", string(spanContextJson))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
