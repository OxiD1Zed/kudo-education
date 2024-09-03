package speller

type SpellerResult struct {
	Code uint     `json:"code"`
	Pos  uint     `json:"pos"`
	Row  uint     `json:"row"`
	Col  uint     `json:"col"`
	Len  uint     `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}
