package voicevox

type Mora struct {
	Text            string  `json:"text"`
	Consonant       string  `json:"consonant,omitempty"`
	ConsonantLength float64 `json:"consonant_length,omitempty"`
	Vowel           string  `json:"vowel"`
	VowelLength     float64 `json:"vowel_length"`
	Pitch           float64 `json:"pitch"`
}

type AccentPhrase struct {
	Moras           []Mora `json:"moras"`
	Accent          int    `json:"accent"`
	PauseMora       *Mora  `json:"pause_mora,omitempty"`
	IsInterrogative bool   `json:"is_interrogative"`
}

type AudioQueryResponse struct {
	AccentPhrases      []AccentPhrase `json:"accent_phrases"`
	SpeedScale         float64        `json:"speedScale"`
	PitchScale         float64        `json:"pitchScale"`
	IntonationScale    float64        `json:"intonationScale"`
	VolumeScale        float64        `json:"volumeScale"`
	PrePhonemeLength   float64        `json:"prePhonemeLength"`
	PostPhonemeLength  float64        `json:"postPhonemeLength"`
	OutputSamplingRate int            `json:"outputSamplingRate"`
	OutputStereo       bool           `json:"outputStereo"`
	Kana               string         `json:"kana"`
}
