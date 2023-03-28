package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
