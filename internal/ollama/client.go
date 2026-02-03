package ollama

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	Model   string
	HTTP    *http.Client
}

type GenerateRequest struct {
	Model  string   `json:"model"`
	Prompt string   `json:"prompt"`
	Images []string `json:"images,omitempty"`
	Stream bool     `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func NewClient(model string) *Client {
	baseURL := os.Getenv("OLLAMA_HOST")
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &Client{
		BaseURL: baseURL,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (c *Client) AnalyzeImage(imagePath, prompt string) (string, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler imagem: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)

	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Images: []string{base64Image},
		Stream: false,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("erro ao serializar request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/api/generate", bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao conectar com Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama retornou status %d", resp.StatusCode)
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return result.Response, nil
}

type ImageAnalysis struct {
	Title       string
	Description string
}

func (c *Client) AnalyzeScreenshot(imagePath string) (*ImageAnalysis, error) {
	prompt := `Esta é uma captura de tela de um software, terminal, IDE ou navegador.
Você é um desenvolvedor analisando um bug ou problema técnico.

Analise a imagem e forneça:
1. Um título curto (máximo 80 caracteres) descrevendo o problema ou erro
2. Uma descrição técnica do que você vê: mensagens de erro, stack traces, problemas de UI, etc.

Se houver texto de erro visível, transcreva-o exatamente.
Se for código, identifique a linguagem e o problema.

Responda em português no formato:
TITULO: <título técnico do problema>
DESCRICAO: <descrição técnica detalhada>`

	response, err := c.AnalyzeImage(imagePath, prompt)
	if err != nil {
		return nil, err
	}

	if os.Getenv("DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, "[DEBUG] Resposta do Ollama:\n%s\n---\n", response)
	}

	return parseAnalysisResponse(response), nil
}

func parseAnalysisResponse(response string) *ImageAnalysis {
	response = strings.TrimSpace(response)

	analysis := &ImageAnalysis{
		Title:       "Análise de screenshot",
		Description: response,
	}

	if response == "" {
		analysis.Description = "(O modelo não retornou uma descrição)"
		return analysis
	}

	lines := strings.Split(response, "\n")
	var descStartIndex int = -1

	for i, line := range lines {
		lineUpper := strings.ToUpper(strings.TrimSpace(line))

		if strings.HasPrefix(lineUpper, "TITULO:") || strings.HasPrefix(lineUpper, "TÍTULO:") || strings.HasPrefix(lineUpper, "TITLE:") {
			colonIdx := strings.Index(line, ":")
			if colonIdx != -1 && colonIdx < len(line)-1 {
				analysis.Title = strings.TrimSpace(line[colonIdx+1:])
			}
		}

		if strings.HasPrefix(lineUpper, "DESCRICAO:") || strings.HasPrefix(lineUpper, "DESCRIÇÃO:") || strings.HasPrefix(lineUpper, "DESCRIPTION:") {
			colonIdx := strings.Index(line, ":")
			descStartIndex = i
			if colonIdx != -1 && colonIdx < len(line)-1 {
				firstPart := strings.TrimSpace(line[colonIdx+1:])
				if firstPart != "" {
					lines[i] = firstPart
				} else {
					descStartIndex = i + 1
				}
			}
			break
		}
	}

	if descStartIndex >= 0 && descStartIndex < len(lines) {
		analysis.Description = strings.TrimSpace(strings.Join(lines[descStartIndex:], "\n"))
	}

	return analysis
}
