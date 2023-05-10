package tokenizer

import (
	"errors"

	tiktokenizer "github.com/tiktoken-go/tokenizer"
)

var _tokenizer tiktokenizer.Codec

func init() {
	enc, err := tiktokenizer.Get(tiktokenizer.Cl100kBase)
	if err != nil {
		panic("error getting tokenizer")
	}
	_tokenizer = enc
}

func CountTokens(content string) (int, error) {
	ids, _, err := _tokenizer.Encode(content)
	if err != nil {
		return 0, errors.New("error couting tokens while encoding")
	}
	return len(ids), err
}

func OffsetTokens(content string, from int, to int) (string, error) {
	ids, _, err := _tokenizer.Encode(content)
	if err != nil {
		return "", errors.New("error couting tokens while encoding")
	}

	start := from
	end := to
	if to > len(ids) {
		end = len(ids)
		println("warning offset out of bounds, truncating to: ", end, "/", to)
	}

	offsetIds := ids[start:end]
	s, _ := _tokenizer.Decode(offsetIds)
	return s, err
}
