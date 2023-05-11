package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/r3labs/sse/v2"
)

const chatApiEndpoint = "https://api.openai.com/v1/chat/completions"

func SendMessage(systemMessage string, userMessage string) string {
	var messages []Message

	if len(systemMessage) > 0 {
		sysMsg := Message{
			Role:    "system",
			Content: systemMessage,
		}
		messages = append(messages, sysMsg)
	}

	msg := Message{
		Role:    "user",
		Content: userMessage,
	}

	messages = append(messages, msg)

	request := OpenaiChatRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
		Stream:   false,
	}

	requestJSON, _ := json.Marshal(request)

	log.Println(string(requestJSON))

	req, err := http.NewRequest("POST", chatApiEndpoint, bytes.NewBuffer(requestJSON))
	if err != nil {
		panic(err)
	}
	apiKey := os.Getenv("OPENAI_API_KEY")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			panic(err)
		}
	}(req.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var response OpenaiChatResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		println("Response Error: ", err.Error())
	}

	returnMessage := response.Choices[0].Message.Content
	return returnMessage
}

type customTransport struct {
	http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Stream(userMessage string) {
	var messages []Message

	msg := Message{
		Role:    "user",
		Content: userMessage,
	}

	messages = append(messages, msg)

	request := OpenaiChatRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
		Stream:   true,
	}

	requestJSON, _ := json.Marshal(request)

	client := &http.Client{
		Transport: &customTransport{
			RoundTripper: http.DefaultTransport,
		},
	}

	sseClient := sse.NewClient(chatApiEndpoint)
	sseClient.Connection = client
	sseClient.Method = "POST"
	sseClient.Body = bytes.NewBuffer([]byte(requestJSON))
	c := make(chan *sse.Event)
	sseClient.SubscribeChanRaw(c)
	speackChanHandler(c)
}

func StreamChannel(messages []Message, ch chan *sse.Event) {
	request := OpenaiChatRequest{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
		Stream:   true,
	}

	requestJSON, _ := json.Marshal(request)

	client := &http.Client{
		Transport: &customTransport{
			RoundTripper: http.DefaultTransport,
		},
	}

	sseClient := sse.NewClient(chatApiEndpoint)
	sseClient.Connection = client
	sseClient.Method = "POST"
	sseClient.Body = bytes.NewBuffer([]byte(requestJSON))
	sseClient.SubscribeChanRaw(ch)
}

func speackHandler(msg *sse.Event) {
	if string(msg.Data) == "[DONE]" {
		return
	}
	var jsonData SseResponse
	err := json.Unmarshal([]byte(msg.Data), &jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s", jsonData.Choices[0].Delta.Content)
}

func speackChanHandler(c chan *sse.Event) {
	for msg := range c {
		if string(msg.Data) == "[DONE]" {
			fmt.Println()
			return
		}
		var jsonData SseResponse
		err := json.Unmarshal([]byte(msg.Data), &jsonData)
		if err != nil {
			fmt.Println(err)
			return
		}

		content := jsonData.Choices[0].Delta.Content
		if strings.Contains(content, "。") {
			fmt.Println(content)
		} else {
			fmt.Printf("%s", content)
		}
	}
}
