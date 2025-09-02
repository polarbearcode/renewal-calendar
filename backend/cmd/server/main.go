package main

import (
    "fmt"
    "log"
    "net/http"

    "myapp/internal/handlers"
	
	"github.com/joho/godotenv"
	 "github.com/unidoc/unipdf/v3/common/license"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	err = license.SetMeteredKey("f053cc479f2fd4c1db4a9c83a2de540196296d306f56539eba6027e15cdb5619")
	if err != nil {
		log.Fatalf("Failed to set UniPDF license: %v", err)
	}

    mux := http.NewServeMux()

    mux.HandleFunc("/upload", handlers.UploadHandler)
	mux.HandleFunc("/parse", handlers.ParseHandler)
	mux.HandleFunc("/calendarData", handlers.CalendarDataHandler)

    fmt.Println("🚀 Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}