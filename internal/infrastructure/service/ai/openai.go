package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAI struct {
	apiKey string
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{apiKey: apiKey}
}

func (o *OpenAI) Ask(input *ProviderInput) (*ProviderOutput, error) {
	requestBody, _ := json.Marshal(input)
	payload := bytes.NewBuffer(requestBody)
	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", payload)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var parsedResponseBody apiResponse
	err = json.Unmarshal(responseBody, &parsedResponseBody)
	if err != nil {
		return nil, err
	}
	text := ""
	if len(parsedResponseBody.Output) > 0 && len(parsedResponseBody.Output[0].Content) > 0 {
		text = parsedResponseBody.Output[0].Content[0].Text
	}
	return &ProviderOutput{Status: parsedResponseBody.Status, Text: text}, nil
}

type apiResponse struct {
	Status string `json:"status"`
	Output []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}
