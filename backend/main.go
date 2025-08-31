package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"


	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client
var bucketName = "vendor-renewals"

func main() {
	// Load default config (needed for endpoint resolver etc.)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-1"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// Initialize S3 client with hardcoded credentials
	s3Client = s3.NewFromConfig(cfg)

	if err != nil {
		log.Fatal("Failed to retrieve credentials:", err)
	}

	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("🚀 Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func enableCORS(w http.ResponseWriter, r *http.Request) bool {
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // or "*" for all origins
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Upload-Timestamp")

    if r.Method == http.MethodOptions {
        // Preflight request, return OK without further handling
        w.WriteHeader(http.StatusNoContent)
        return false
    }
    return true
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	 if !enableCORS(w, r) {
        return
    }


	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploadDateTime := r.Header.Get("X-Upload-Timestamp")

	fmt.Println("Upload DateTime:", uploadDateTime)

	// Parse uploaded file (max 10MB)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too big", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Could not get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file into memory
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Upload to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uploadDateTime + "/" + header.Filename),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    "private",
	})
	if err != nil {
		http.Error(w, "Failed to upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "✅ Uploaded %s to S3\n", header.Filename)
}
