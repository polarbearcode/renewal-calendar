package services




import (
	
	"fmt"
	"bytes"
    "log"
	"os"
	"context"


	"encoding/json"
	 "io/ioutil"
	 "net/http"
	 "mime/multipart"

	
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"

	"myapp/ai"

)



func ParseFiles(w http.ResponseWriter, files []*multipart.FileHeader) string {
    // TODO: read from S3, parse content, return result

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	aiClient := ai.NewOpenRouterClient(apiKey)

	for _, header := range files {
		fmt.Println("Parsing:", header.Filename)
		file, err := header.Open()
		if err != nil {
			http.Error(w, "Failed to open file: "+err.Error(), http.StatusInternalServerError)
			return ""
		}

		fileBytes, err := ioutil.ReadAll(file)
		file.Close()
		if err != nil {
			http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
			return ""
		}

		text, err := ExtractTextFromPDFBytes(fileBytes)
		if err != nil {
			log.Printf("failed to extract text from PDF: %v", err)
			continue
		}

		result, err := aiClient.Chat(context.Background(), text)
		if err != nil {
			log.Printf("failed to chat with AI: %v", err)
			continue
		}

		fmt.Println(result)
	}

	return ""
}

// ExtractTextFromPDFBytes extracts all text from a PDF provided as []byte.
func ExtractTextFromPDFBytes(pdfBytes []byte) (string, error) {

	reader := bytes.NewReader(pdfBytes)

	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return "", fmt.Errorf("failed to create PDF reader: %w", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("failed to get number of pages: %w", err)
	}

	var fullText string
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			log.Printf("Warning: failed to get page %d: %v", i, err)
			continue
		}

		ex, err := extractor.New(page)
		if err != nil {
			log.Printf("Warning: failed to create extractor for page %d: %v", i, err)
			continue
		}

		text, err := ex.ExtractText()
		if err != nil {
			log.Printf("Warning: failed to extract text from page %d: %v", i, err)
			continue
		}

		fullText += text + "\n"
	}

	return fullText, nil
}

func SendTextToOpenAI(text string) string {
    apiKey := "sk-or-5af618b4e724a34c038585d197e8303dbd5a24e62688e3a8800efad838abd285" // Replace with your key
    url := "https://openrouter.ai/api/v1/chat/completions"

    // Build request body
    body := map[string]interface{}{
        "model": "openai/gpt-4o", // or "gpt-4o-mini"
        "messages": []map[string]string{
            {"role": "user", "content": text},
        },
    }

    bodyBytes, err := json.Marshal(body)
    if err != nil {
        log.Fatalf("Failed to marshal body: %v", err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
    if err != nil {
        log.Fatalf("Failed to create request: %v", err)
    }

	fmt.Println(apiKey)
    req.Header.Set("Authorization", "Bearer "+apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
    Transport: http.DefaultTransport, // ensures no AWS signing
}


    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Request failed: %v", err)
    }
    defer resp.Body.Close()

    respBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response: %v", err)
    }

    // Print raw response for debugging
    fmt.Println("Raw response:", string(respBytes))

    // Parse JSON response
    var parsedResp struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }

    if err := json.Unmarshal(respBytes, &parsedResp); err != nil {
        log.Fatalf("Failed to parse JSON: %v", err)
    }

    if len(parsedResp.Choices) == 0 {
        return ""
    }

    return parsedResp.Choices[0].Message.Content
}