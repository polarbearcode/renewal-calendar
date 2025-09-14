package handlers

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
	"log"
	"regexp"
	"time"
	"strconv"
	"context"
	"bytes"

    "github.com/ledongthuc/pdf"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/time/rate"
)

// Struct for incoming JSON payload
type UploadPayload struct {
    FileBytesMap map[string]string `json:"fileBytesMap"` // Base64 strings
}

func extractTextFromPDFBytes(fileBytes []byte) (string, error) {
    reader := bytes.NewReader(fileBytes)

    // Open PDF from memory instead of file
    pdfReader, err := pdf.NewReader(reader, int64(len(fileBytes)))
    if err != nil {
        return "", err
    }

    var text strings.Builder
    numPages := pdfReader.NumPage()
    for i := 1; i <= numPages; i++ {
        page := pdfReader.Page(i)
        pageText, err := page.GetPlainText(nil)
        if err != nil {
            fmt.Println("Failed to extract text from page", i, err)
            continue
        }
        text.WriteString(pageText)
    }

    return text.String(), nil
}

type FileResult struct {
    Filename string `json:"filename"`
    Success  bool   `json:"success"`
    Data     string `json:"data,omitempty"`
    Error    string `json:"error,omitempty"`
}

func ParseHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ParseHandler reached! Method=%s Path=%s", r.Method, r.URL.Path)
	if !enableCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload UploadPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		http.Error(w, "OPENAI_API_KEY not set", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	// Limit to 1 request every 2 seconds
	limiter := rate.NewLimiter(rate.Every(2*time.Second), 1)

	results := make(map[string]map[string]interface{})

	

	for filename, b64 := range payload.FileBytesMap {
		fileBytes, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			fmt.Println("Failed to decode Base64 for", filename, err)
			continue
		}

		text, err := extractTextFromPDFBytes(fileBytes)
		if err != nil {
			fmt.Println("Failed to extract text from", filename, err)
			continue
		}

		// Wait for rate limiter before sending API request
		if err := limiter.Wait(context.Background()); err != nil {
			fmt.Println("Rate limiter error:", err)
			continue
		}

		chatReq := map[string]interface{}{
			"model": "deepseek/deepseek-r1:free",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": fmt.Sprintf("Give me the seller name, effective date, renewal date, and whether this contract auto renews in a very parsable string for this file:\n---\n%s", text),
				},
			},
		}

		bodyBytes, _ := json.Marshal(chatReq)
		req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", strings.NewReader(string(bodyBytes)))
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		respBody, err := callOpenRouterWithRetry(req, 3, client)
		if err != nil {
			fmt.Println("API call failed for", filename, err)
			continue
		}

		fmt.Println("Here", respBody)

		responseMap, err := responseToMap(respBody)
		
		if err != nil {
			fmt.Println("Failed to parse response for", filename, err)
			continue
		}

		results[filename] = responseMap

		// Optionally send each file’s parsed result to the database
		sendMapToDB(responseMap)
		fmt.Println("Processed file:", filename)
	}

	// Return combined results for all files
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func responseToMap(respBody []byte) (map[string]interface{}, error) {

	type Response struct {
    	choices string `json:"choices"`
    	// Add other fields as needed
	}

	type ChatResponse struct {
		Choices []struct {
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	fmt.Println("Here: ", respBody)

	var parsedResp ChatResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		fmt.Println("here")
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Println("Parsed response:", parsedResp.Choices[0].Message.Content)
	fmt.Println("Woo")

	re := regexp.MustCompile(`(?i)(?:name|seller_name|seller name):\s*(.+?)\s*(?:,|;|\||\n)\s*(?:effective date|effective_date|Effective Date):\s*(.+?)\s*(?:,|;|\||\n)\s*(?:renewal[-_ ]?date):\s*(.+?)\s*(?:,|;|\||\n)\s*(?:autorenew|auto_renew|auto_renews):\s*(\w+)`)

    matches := re.FindStringSubmatch(parsedResp.Choices[0].Message.Content)

	fmt.Println("Parsed response:", matches)

	if len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse response")
	}

	return map[string]interface{}{
		"name":          matches[1],
		"effective_date": matches[2],
		"renewal_date":  matches[3],
		"autorenew":     matches[4],
	}, nil
}

func sendMapToDB(data map[string]interface{}) error {

	type Contract struct {
		Name          string `json:"seller"`
		EffectiveDate string `json:"effective_date"`
		RenewalDate   string `json:"renewal_date"`
		Autorenew     bool `json:"autorenew"`
	}

	contract := Contract{
		Name:          data["name"].(string),
		EffectiveDate: data["effective_date"].(string),
		RenewalDate:   data["renewal_date"].(string),
	}

	
	contract_autorenew := data["autorenew"]

	if strings.Contains(strings.ToLower(contract_autorenew.(string)), "yes") {
		contract.Autorenew = true
	} else {
		contract.Autorenew = false
	}

	apiURL := os.Getenv("API_URL")
	apiKey := os.Getenv("API_KEY")
	client, err := supabase.NewClient(apiURL, apiKey, &supabase.ClientOptions{})
	if err != nil {
		return fmt.Errorf("failed to create Supabase client: %w", err)
	}

	

	res, count, err := client.From("contracts").Insert(
		[]Contract{contract}, // must still be a slice
		false,                // upsert?
		"",                   // onConflict
		"representation",     // returning ("representation" or "minimal")
		"",                   // count
	).Execute()

	if err != nil {
		log.Fatalf("failed to insert: %v", err)
	}

	fmt.Println(count, string(res)) // response from Supabase

	return nil
}

func callOpenRouterWithRetry(req *http.Request, maxRetries int, client *http.Client) ([]byte, error) {
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		// Retry if rate-limited
		if resp.StatusCode == http.StatusTooManyRequests { // 429
			retryAfter := 5 * time.Second // default
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if sec, parseErr := strconv.Atoi(ra); parseErr == nil {
					retryAfter = time.Duration(sec) * time.Second
				}
			}
			fmt.Printf("429 received, retrying after %v...\n", retryAfter)
			time.Sleep(retryAfter)
			continue
		}

		// If other errors, return with body for debugging
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("OpenRouter API error %d: %s", resp.StatusCode, string(bodyBytes))
		}

		// Successful response
		return bodyBytes, nil
	}

	return nil, fmt.Errorf("exceeded max retries due to rate limiting")
}
