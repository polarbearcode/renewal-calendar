package handlers

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    "github.com/ledongthuc/pdf"
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

	fmt.Println(payload)

    var allText strings.Builder
    for filename, b64 := range payload.FileBytesMap {
        fileBytes, err := base64.StdEncoding.DecodeString(b64)
		fmt.Println(66)
        if err != nil {
            fmt.Println("Failed to decode Base64 for", filename, err)
            continue
        }

		fmt.Println(72)

        text, err := extractTextFromPDFBytes(fileBytes)
        if err != nil {
            fmt.Println("Failed to extract text from", filename, err)
            continue
        }

		fmt.Println(80)

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
            {"role": "user", "content": allText.String()},
        },
    }



    bodyBytes, _ := json.Marshal(chatReq)
    req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", strings.NewReader(string(bodyBytes)))
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

	fmt.Println(108)

    client := &http.Client{}

	fmt.Println(112)

	
    resp, err := client.Do(req)
    if err != nil {
		fmt.Println(err)
        http.Error(w, "Error calling OpenAI API: "+err.Error(), http.StatusInternalServerError)
        return
    }

    defer resp.Body.Close()
    respBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("OpenAI response:", string(respBody))

    w.Header().Set("Content-Type", "application/json")
    w.Write(respBody)
}