package handlers 

import (
	"net/http"
	"fmt"
	"github.com/supabase-community/supabase-go"
	"os"

)

func CalendarDataHandler(w http.ResponseWriter, r *http.Request) {

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
		fmt.Println("cannot initalize client", err)
	}

	data, count, err := client.From("contracts").Select("*", "exact", false).Execute()

	if err != nil {
		fmt.Println("failed to fetch data:", count, err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}
