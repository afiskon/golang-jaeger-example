package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"go.opencensus.io/trace"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/afiskon/golang-jaeger-example"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	spanContextJson := r.Header.Get("X-Span-Context")
	var spanContext trace.SpanContext
	err := json.Unmarshal([]byte(spanContextJson), &spanContext)
	if err != nil {
		logrus.Errorf("PostHandler, json.Unmarshal: %v\n", err)
		return
	}

	_ /* ctx2 */, span := trace.StartSpanWithRemoteParent(
		context.Background(), "PostHandler", spanContext)
	defer span.End()

	logrus.Debugf("PostHandler called, traceId = %s\n", span.SpanContext().TraceID.String())

	h := w.Header()
	h.Set("Content-Type", "text/html")
	w.WriteHeader(200)
	_, err = w.Write([]byte("<h1>PONG</h1>\n"))
	if err != nil {
		logrus.Errorf("PostHandler, w.Write: %v\n", err)
	}
}

var listenAddr string

func main() {
	flag.StringVar(&listenAddr, "listen", ":8080",
		"Listen address")
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)
	utils.InitJaeger("serer")

	logrus.Infof("Starting HTTP server at %s...",
		listenAddr)

	r := mux.NewRouter().
		PathPrefix("/api/v1").
		Path("/ping").
		Subrouter()

	r.Methods("POST").
		HandlerFunc(PostHandler)

	http.Handle("/", r)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		logrus.Fatalf("http.ListenAndServe: %v\n", err)
	}

	logrus.Info("HTTP server terminated\n")
}
