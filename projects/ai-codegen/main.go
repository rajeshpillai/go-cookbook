package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"

	"go-cookbook/projects/ai-codegen/mcp"
)

// StringArray is a custom type that can unmarshal either a JSON array of strings
// or an object (in which case it extracts the keys as a slice of strings).
type StringArray []string

func (sa *StringArray) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a slice of strings.
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*sa = arr
		return nil
	}

	// If not, try unmarshaling as an object and collect its keys.
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}
	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}
	*sa = keys
	return nil
}

// CodeResponse now uses StringArray for FileStructure.
type CodeResponse struct {
	FileStructure StringArray       `json:"fileStructure"`
	CodeFiles     map[string]string `json:"codeFiles"`
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found; using environment variables.")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		color.Red("⚠️ Missing OpenAI API Key. Set OPENAI_API_KEY in env.")
		os.Exit(1)
	}

	// Clear terminal
	fmt.Print("\033[2J\033[H")
	color.Cyan("\n🚀 AI Code Generator CLI\n")

	appType, err := selectAppType()
	if err != nil {
		log.Fatalf("Error selecting app type: %v", err)
	}

	var mcpContent string
	switch appType {
	case "simple-code":
		mcpContent = mcp.SimpleCode
	case "term-app":
		mcpContent = mcp.TermApp
	case "web-app":
		mcpContent = mcp.WebApp
	default:
		color.Red("Invalid app type")
		os.Exit(1)
	}

	model, err := selectModel()
	if err != nil {
		log.Fatalf("Error selecting model: %v", err)
	}

	projectName, err := prompt("Enter a name for your project folder:", "my-app")
	if err != nil {
		log.Fatalf("Error reading project name: %v", err)
	}

	projectIdea, err := getInputWithEditor("")
	if err != nil {
		log.Fatalf("Error reading project idea: %v", err)
	}

	confirmedPrompt, err := previewPrompt(projectIdea)
	if err != nil {
		log.Fatalf("Error during prompt preview: %v", err)
	}

	// Generate code using streaming
	rawResponse, err := generateCodeWithModelStream(confirmedPrompt, mcpContent, model, apiKey)
	if err != nil {
		log.Fatalf("Error generating code: %v", err)
	}

	codeResp, err := parseJsonResponse(rawResponse)
	if err != nil {
		log.Fatalf("Error parsing code response: %v", err)
	}

	err = saveFiles(codeResp, projectName)
	if err != nil {
		log.Fatalf("Error saving files: %v", err)
	}

	color.Green("\n🎉 Project '%s' created successfully!\n", projectName)
}

func selectAppType() (string, error) {
	var appType string
	prompt := &survey.Select{
		Message: "Which type of app do you want to generate?",
		Options: []string{"simple-code", "term-app", "web-app"},
		Default: "simple-code",
	}
	err := survey.AskOne(prompt, &appType)
	return appType, err
}

func selectModel() (string, error) {
	var model string
	prompt := &survey.Select{
		Message: "Which model would you like to use?",
		Options: []string{"gpt-4", "codellama:7b-instruct", "wizardcoder:7b"},
		Default: "gpt-4",
	}
	err := survey.AskOne(prompt, &model)
	return model, err
}

func prompt(message, defaultValue string) (string, error) {
	var response string
	p := &survey.Input{
		Message: message,
		Default: defaultValue,
	}
	err := survey.AskOne(p, &response)
	return response, err
}

func getInputWithEditor(defaultText string) (string, error) {
	var response string
	p := &survey.Editor{
		Message:  "Describe your app idea (multi-line):",
		FileName: "*.txt",
		Default:  defaultText,
	}
	err := survey.AskOne(p, &response)
	return strings.TrimSpace(response), err
}

func previewPrompt(projectIdea string) (string, error) {
	for {
		color.Cyan("\n📄 Here's what you wrote:\n")
		fmt.Println(projectIdea)
		var choice string
		p := &survey.Select{
			Message: "Proceed with this input?",
			Options: []string{"yes", "edit", "cancel"},
			Default: "yes",
		}
		err := survey.AskOne(p, &choice)
		if err != nil {
			return "", err
		}
		if choice == "yes" {
			return projectIdea, nil
		} else if choice == "cancel" {
			color.Red("❌ Cancelled.")
			os.Exit(0)
		} else if choice == "edit" {
			newIdea, err := getInputWithEditor(projectIdea)
			if err != nil {
				return "", err
			}
			projectIdea = newIdea
		}
	}
}

func generateCodeWithModelStream(projectIdea, mcpContent, model, apiKey string) (string, error) {
	color.Blue("\n💡 Generating code with %s...\n", model)
	// Combine system prompt (MCP) and user input.
	fullPrompt := mcpContent + "\n" + projectIdea

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	var fullResponse string
	var err error
	if strings.Contains(model, ":") {
		fullResponse, err = generateOllamaStream(fullPrompt, model)
	} else {
		fullResponse, err = generateOpenAIStream(fullPrompt, model, apiKey)
	}
	if err != nil {
		return "", err
	}
	fmt.Println()
	return fullResponse, nil
}

func generateOpenAIStream(prompt, model, apiKey string) (string, error) {
	client := openai.NewClient(apiKey)
	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: prompt[:strings.Index(prompt, "\n")]},
			{Role: "user", Content: prompt},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var result strings.Builder
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		// Use HiBlack color to mimic gray text.
		color.New(color.FgHiBlack).Printf("%s", response.Choices[0].Delta.Content)
		result.WriteString(response.Choices[0].Delta.Content)
	}
	return result.String(), nil
}

func generateOllamaStream(prompt, model string) (string, error) {
	baseURL := os.Getenv("OLLAMA_API_BASE_URL")
	if baseURL == "" {
		return "", errors.New("OLLAMA_API_BASE_URL not set")
	}
	payload := map[string]interface{}{
		"model":  model,
		"stream": true,
		"messages": []map[string]string{
			{"role": "system", "content": prompt[:strings.Index(prompt, "\n")]},
			{"role": "user", "content": prompt},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", baseURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		color.New(color.FgHiBlack).Printf("%s", line)
		result.WriteString(line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return result.String(), nil
}

func extractJsonFromResponse(text string) (string, error) {
	re := regexp.MustCompile("(?s)```json(.*?)```")
	matches := re.FindStringSubmatch(text)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1]), nil
	}
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start != -1 && end != -1 && end > start {
		return strings.TrimSpace(text[start : end+1]), nil
	}
	return "", errors.New("no JSON found in the AI response")
}

func parseJsonResponse(raw string) (*CodeResponse, error) {
	jsonStr, err := extractJsonFromResponse(raw)
	if err != nil {
		return nil, err
	}
	var codeResp CodeResponse
	err = json.Unmarshal([]byte(jsonStr), &codeResp)
	if err != nil {
		return nil, err
	}
	if codeResp.CodeFiles == nil || len(codeResp.CodeFiles) == 0 {
		return nil, errors.New("missing or invalid 'codeFiles' in parsed response")
	}
	return &codeResp, nil
}

func saveFiles(codeResp *CodeResponse, projectName string) error {
	baseDir := filepath.Join("output", projectName)
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return err
	}
	color.Green("\n📁 Writing files to ./%s/\n", projectName)
	for _, filePath := range codeResp.FileStructure {
		content, ok := codeResp.CodeFiles[filePath]
		if !ok {
			content = "// No content provided"
		}
		fullPath := filepath.Join(baseDir, filePath)
		os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		err := ioutil.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			return err
		}
		color.Yellow("✅ %s", filePath)
	}
	return nil
}
