package text_validator

import (
	"fmt"
	"kode-education/pkg/speller"
	"regexp"
)

type TextValidator struct {
	speller speller.YandexSpeller
}

func New(speller speller.YandexSpeller) *TextValidator {
	return &TextValidator{
		speller: speller,
	}
}

func (v *TextValidator) ValidateTextMultiLang(text string) (string, error) {
	const (
		op        = "TextValidator.ValidateTextMultiLang"
		patternEN = ".*[A-Za-z]+.*"
		patternRU = ".*[А-Яа-я]+.*"
		patternUK = ".*[А-ЩЬЮЯҐЄІЇа-щьюяґєії]+.*"
	)
	arrayPatterns := map[string]string{
		speller.LangEN: patternEN,
		speller.LangRU: patternRU,
	}

	for lang, pattern := range arrayPatterns {
		if has, _ := regexp.MatchString(pattern, text); has {
			var err error
			text, err = v.ValidateText(text, lang)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}
		}
	}

	return text, nil
}

func (v *TextValidator) ValidateText(text, lang string) (string, error) {
	const op = "TextValidator.ValidateText"

	spellRes, err := v.speller.CheckText(text, lang, speller.OptionDefault)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return v.replaceIncorrectWord(text, spellRes), nil
}

func (v *TextValidator) replaceIncorrectWord(text string, spellRes []speller.SpellerResult) string {
	diff := 0
	runeText := []rune(text)
	for i := 0; i < len(spellRes); i++ {
		if len(spellRes[i].S) < 1 {
			continue
		}

		startWord := int(spellRes[i].Pos) + diff
		endWord := startWord + int(spellRes[i].Len)
		correctWord := []rune(spellRes[i].S[0])

		runeText = append(runeText[:startWord], append(correctWord, runeText[endWord:]...)...)

		diff += len(correctWord) - int(spellRes[i].Len)
	}
	return string(runeText)
}
