package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func main() {
	apiKey, err := loadApiKey()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) <= 1 {
		printHelpManual()
		os.Exit(0)
	}

	var country string
	var arg string

	if len(os.Args) > 2 {
		country = os.Args[1]
		arg = os.Args[2]
	} else {
		arg = os.Args[1]
	}

	switch arg {
	case "--help", "-h":
		printHelpManual()
	case "--version", "-v":
		printLastVersion()
	default:
		translation, err := postTranslation(arg, country, apiKey)
		if err != nil {
			log.Error("Error translating word")
			os.Exit(1)
		}

		log.Info(fmt.Sprintf("Translation: %s", renderElSelected(translation)))
	}

	os.Exit(0)
}

func printHelpManual() {
	cmds := []string{
		"trad",
		"trad [--help, -h]",
	}
	descs := []string{
		"Run the program",
		"Show this help message",
	}

	fmt.Println("\nUsage: trad [options]")
	for i, cmd := range cmds {
		fmt.Printf("  %-20s %s\n", cmd, descs[i])
	}
}

func getLatestVersion() string {
	return "v0.0.1"
}

func printLastVersion() {
	fmt.Printf("Latest version: %s\n", getLatestVersion())
}

func loadApiKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	return os.Getenv("API_KEY"), nil
}

func postTranslation(word string, country string, apiKey string) (translation string, err error) {
	client := &http.Client{}
	requestBody := map[string]any{
		"text": []string{
			word,
		},
	}

	if country != "" {
		requestBody["target_lang"] = strings.ToUpper(country)
	} else {
		requestBody["target_lang"] = "FR"
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api-free.deepl.com/v2/translate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var deeplResp struct {
		Translations []struct {
			DetectedSourceLanguage string `json:"detected_source_language"`
			Text                   string `json:"text"`
		} `json:"translations"`
	}

	err = json.NewDecoder(resp.Body).Decode(&deeplResp)
	if err != nil {
		return "", err
	}

	if len(deeplResp.Translations) > 0 {
		return deeplResp.Translations[0].Text, nil
	}

	return "", nil
}

func renderElSelected(el string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", "38;2;214;112;214", el)
}
