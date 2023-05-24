package tokenizer

import (
	"errors"
	"strings"

	tiktokenizer "github.com/tiktoken-go/tokenizer"
	"github.com/tzapio/tzap/internal/logging/tl"
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

func OffsetTokens(content string, from int, to int) (string, int, error) {
	ids, strs, err := _tokenizer.Encode(content)
	if err != nil {
		return "", 0, errors.New("error couting tokens while encoding")
	}

	start := from
	end := to
	if to > len(ids) {
		end = len(ids)
		tl.Logger.Println("warning offset out of bounds, truncating to: ", end, "/", to)
	}
	offsetStrs := strs[start:end]
	s := strings.Join(offsetStrs, "")
	return s, len(offsetStrs), err
}

func RawTokens(content string) ([]string, error) {
	return []string{}, nil
}
