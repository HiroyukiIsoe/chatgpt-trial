package voicevox

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func Speak(text string, outPath string) {
	audioQuery := getAudioQuery(text)

	b := synthesis(audioQuery)

	out := outPath
	if len(out) == 0 {
		now := time.Now()
		out = "tmp/audio" + now.Format("_20060102150405") + ".wav"
	}

	if err := ioutil.WriteFile(out, b, 0644); err != nil {
		panic(err)
	}
}

func getAudioQuery(text string) AudioQueryResponse {
	apiHost := os.Getenv("VOICEVOX_API_HOST")
	req, err := http.NewRequest("POST", apiHost+"/audio_query", nil)
	if err != nil {
		panic(err)
	}

	query := req.URL.Query()
	query.Add("speaker", "1")
	query.Add("text", text)
	req.URL.RawQuery = query.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var response AudioQueryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		println("Voicevox Response Error: ", err.Error())
	}

	return response
}

func synthesis(query AudioQueryResponse) []byte {
	apiHost := os.Getenv("VOICEVOX_API_HOST")
	jsonBody, _ := json.Marshal(query)

	req, err := http.NewRequest("POST", apiHost+"/synthesis", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}

	httpQuery := req.URL.Query()
	httpQuery.Add("speaker", "1")
	req.URL.RawQuery = httpQuery.Encode()

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "audio/wav")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return body
}
