package main

import (
	"chatgpt-trial/pkg/openai"
	"chatgpt-trial/pkg/voicevox"
	"log"
	"strings"
)

func main() {
	var sb strings.Builder
	sb.WriteString("文章を書くときは次のルールを忠実に守ってください。")
	sb.WriteString("語尾に「なのだ」もしくは「のだ」をつけること。")
	sb.WriteString("小学校6年生に説明するようにわかりやすく説明すること。")
	sb.WriteString("可能な限り簡潔にな文章にすること。")
	sb.WriteString("一人称は「ぼく」とすること。")

	question := "チキン南蛮について教えて下さい"

	log.Println("チキン南蛮について聞いてみます。")
	responseMessage := openai.SendMessage(sb.String(), question)

	log.Println("ChatGPTに回答: ", responseMessage)

	log.Println("音声を生成します")
	voicevox.Speak(responseMessage, "")
}
