package tokenizer

import (
	"errors"
	"strings"

	"github.com/tzapio/tokenizer/codec"
	enc "github.com/tzapio/tokenizer/codec/cl100k_base"
	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/util/singlewait"
)

type Tokenizer struct {
	tokenizer *singlewait.SingleWait[*codec.Codec]
}

func NewTokenizer() *Tokenizer {
	waitTokenizer := singlewait.New(func() *codec.Codec {
		tl.Logger.Println("Initiating Tokenizer Client")
		tokenizer := enc.NewCl100kBase()
		tl.Logger.Println("Done - Initializing Tokenizer Client")
		return tokenizer
	})
	t := &Tokenizer{
		tokenizer: waitTokenizer,
	}

	return t
}
func (t *Tokenizer) CountTokens(content string) (int, error) {
	ids, _, err := t.tokenizer.GetData().Encode(content)
	if err != nil {
		return 0, errors.New("error couting tokens while encoding")
	}
	return len(ids), err
}

func (t *Tokenizer) OffsetTokens(content string, from int, to int) (string, int, error) {
	ids, strs, err := t.tokenizer.GetData().Encode(content)
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

func (t *Tokenizer) RawTokens(content string) ([]string, error) {
	return []string{}, nil
}
