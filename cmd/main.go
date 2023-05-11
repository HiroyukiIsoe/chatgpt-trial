package main

import (
	"chatgpt-trial/pkg/openai"
	"chatgpt-trial/pkg/voicevox"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r3labs/sse/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	var messages []openai.Message
	var sb strings.Builder
	sb.WriteString("文章を書くときは次のルールを忠実に守ってください。")
	sb.WriteString("語尾に「なのだ」もしくは「のだ」をつけること。")
	sb.WriteString("小学校6年生に説明するようにわかりやすく説明すること。")
	sb.WriteString("可能な限り簡潔にな文章にすること。")
	sb.WriteString("一人称は「ぼく」とすること。")

	sysMsg := openai.Message{
		Role:    "system",
		Content: sb.String(),
	}

	messages = append(messages, sysMsg)

	question := openai.Message{
		Role:    "user",
		Content: "チキン南蛮について教えて下さい",
	}

	messages = append(messages, question)

	// log.Println("チキン南蛮について聞いてみます。")
	// responseMessage := openai.SendMessage(sb.String(), question)

	// log.Println("ChatGPTの回答: ", responseMessage)

	// log.Println("音声を生成します")
	// voicevox.Speak(responseMessage, "")
	eg, _ := errgroup.WithContext(context.Background())
	ch := make(chan *sse.Event)
	openai.StreamChannel(messages, ch)
	makeVoice(eg, ch)
	if err := eg.Wait(); err != nil {
		panic(err)
	}
}

func makeVoice(eg *errgroup.Group, ch chan *sse.Event) {

	var voiceBlock string
	for msg := range ch {
		if string(msg.Data) == "[DONE]" {
			fmt.Println()
			return
		}
		var jsonData openai.SseResponse
		err := json.Unmarshal([]byte(msg.Data), &jsonData)
		if err != nil {
			fmt.Println(err)
			return
		}

		content := jsonData.Choices[0].Delta.Content
		fmt.Print(content)
		voiceBlock += content
		if strings.Contains(content, "。") {
			fmt.Println()
			voice := voiceBlock
			voiceBlock = ""
			eg.Go(func() error {
				if err := voicevox.Speak(voice, ""); err != nil {
					return err
				}
				return nil
			})
		}
	}
}
