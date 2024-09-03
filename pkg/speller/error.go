package speller

import "errors"

var (
	ErrorTextOverflow        = errors.New("the text is too big")
	ErrorServiceNotAvailable = errors.New("the service is not currently available")
	ErrorInvalidParameters   = errors.New("invalid parameters")
)
