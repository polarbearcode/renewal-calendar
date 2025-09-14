package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/aws/aws-lambda-go/lambda"
    "github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"myapp/internal/handlers"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("Lambda reached! Method=%s Path=%s", r.Method, r.URL.Path)
    fmt.Fprintf(w, `{"status":"ok"}`)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", handlers.ParseHandler)

    adapter := httpadapter.NewV2(mux) // <-- HTTP API v2 adapter
    lambda.Start(adapter.ProxyWithContext)
}