package handlers 

import (
	"net/http"
	"fmt"
	"github.com/supabase-community/supabase-go"
	"os"
	"log"

)

func CalendarDataHandler(w http.ResponseWriter, r *http.Request) {

	    log.Printf("Received request: method=%s path=%s", r.Method, r.URL.Path)


	if !enableCORS(w, r) {
		return
	}

    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	API_URL := os.Getenv("API_URL")
	API_KEY := os.Getenv("API_KEY")

    client, err := supabase.NewClient(API_URL, API_KEY, &supabase.ClientOptions{})

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to initialize Supabase client: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		log.Println(errorMessage)
		return
	}

	data, _, err := client.From("contracts").Select("*", "exact", false).Execute()

	if err != nil {
		errorMessage := fmt.Sprintf("Failed to fetch data from Supabase: %v", err)
		http.Error(w, errorMessage, http.StatusInternalServerError)
		log.Println(errorMessage)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}
