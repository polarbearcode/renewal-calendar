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

    "github.com/ledongthuc/pdf"
	"github.com/supabase-community/supabase-go"
)

// Struct for incoming JSON payload
type UploadPayload struct {
    FileBytesMap map[string]string `json:"fileBytesMap"` // Base64 strings
}

func extractTextFromPDFBytes(fileBytes []byte) (string, error) {
    tmpFile := "temp.pdf"
    if err := ioutil.WriteFile(tmpFile, fileBytes, 0644); err != nil {
        return "", err
    }

    _, pdfReader, err := pdf.Open(tmpFile)
    if err != nil {
        return "", err
    }

    var text strings.Builder
    numPages := pdfReader.NumPage()
    for i := 1; i <= numPages; i++ { // Pages are 1-indexed
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

func ParseHandler(w http.ResponseWriter, r *http.Request) {

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
    var allText strings.Builder
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


        allText.WriteString("\n---\n")
        allText.WriteString(fmt.Sprintf("File: %s\n", filename))
        allText.WriteString(text)
    }



    // Now send allText.String() to OpenAI Chat API
    apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
    if apiKey == "" {
        http.Error(w, "OPENAI_API_KEY not set", http.StatusInternalServerError)
        return
    }

	

    // Prepare OpenRouter Chat API request
    chatReq := map[string]interface{}{
        "model": "deepseek/deepseek-r1:free", // or whichever model you want
        "messages": []map[string]string{
            {"role": "user", "content": "Give me the seller name, effective date, renewal date, and whether this contract auto renews in a very parsable string like name: ..., effective date: ... renewal-date: .... autorenew: yes/no that I can easily parse with a regex" + allText.String()},
        },
    }



    bodyBytes, _ := json.Marshal(chatReq)
    req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", strings.NewReader(string(bodyBytes)))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

	

    client := &http.Client{}

	

	
    resp, err := client.Do(req)
    if err != nil {
		fmt.Println(err)
        http.Error(w, "Error calling OpenAI API: "+err.Error(), http.StatusInternalServerError)
        return
    }

    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)

	responseMap, _ := responseToMap(respBody)

	sendMapToDB(responseMap)

	fmt.Println("Response sent to database:", responseMap)

    w.Header().Set("Content-Type", "application/json")
    w.Write(respBody)
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

	var parsedResp ChatResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		fmt.Println("here")
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Println("Parsed response:", parsedResp.Choices[0].Message.Content)
	fmt.Println("Woo")

	re := regexp.MustCompile(`(?s)name:\s*(.+?)[,\n]\s*effective date:\s*(.+?)[,\n]\s*renewal-date:\s*(.+?)[,\n]\s*autorenew:\s*(\w+)`)

    matches := re.FindStringSubmatch(parsedResp.Choices[0].Message.Content)

	fmt.Println("Parsed response:", matches)

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
