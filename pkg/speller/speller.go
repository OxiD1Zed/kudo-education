package speller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	host  = "https://speller.yandex.net"
	route = "services/spellservice.json"
)

const (
	methodCheckText  = "checkText"
	methodCheckTexts = "checkTexts"
)

const (
	OptionDefault               = 0
	OptionIgnoreDigits          = 2
	OptionIgnoreUrld            = 4
	OptionRepeatWords           = 8
	OptionIgnoreCapitialization = 512
)

const (
	LangRU = "ru"
	LangEN = "en"
	LangUK = "uk"
)

type YandexSpeller struct {
	client http.Client
}

func New(client http.Client) *YandexSpeller {
	return &YandexSpeller{client: client}
}

func (s *YandexSpeller) check(text []string, lang, method string, options int) ([]byte, error) {
	resp, err := s.queryToSpeller(text, lang, method, options)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrorServiceNotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrorInvalidParameters
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorServiceNotAvailable
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *YandexSpeller) queryToSpeller(text []string, lang, method string, options int) (*http.Response, error) {
	const (
		limitText = 10000
	)

	if len(text) > limitText {
		return nil, ErrorTextOverflow
	}

	uri, err := url.JoinPath(host, route, method)
	if err != nil {
		return nil, err
	}

	data := url.Values{
		"text":    text,
		"lang":    {lang},
		"options": {fmt.Sprint(options)},
	}

	return s.client.PostForm(uri, data)
}

func (s *YandexSpeller) CheckText(text, lang string, options int) ([]SpellerResult, error) {
	const op = "YandexSpeller.CheckText"

	body, err := s.check([]string{text}, lang, methodCheckText, options)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.jsonCheckTextToSpellerRes(body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *YandexSpeller) jsonCheckTextToSpellerRes(body []byte) ([]SpellerResult, error) {
	var responce []SpellerResult
	if err := json.Unmarshal(body, &responce); err != nil {
		return nil, err
	}

	return responce, nil
}

func (s *YandexSpeller) CheckTexts(text []string, lang string, options int) ([][]SpellerResult, error) {
	const op = "YandexSpeller.CheckTexts"

	body, err := s.check(text, lang, methodCheckText, options)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.jsonCheckTextsToSpellerRes(body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *YandexSpeller) jsonCheckTextsToSpellerRes(body []byte) ([][]SpellerResult, error) {
	var responce [][]SpellerResult
	if err := json.Unmarshal(body, &responce); err != nil {
		return nil, err
	}

	return responce, nil
}
