package idl_conv

import (
	"errors"
	"strconv"
)

type CharIterator struct {
	i     int
	runes []rune
	n     int
}

func NewCharIterator(str string) *CharIterator {

	runes := []rune(str)

	return &CharIterator{
		i:     0,
		runes: runes,
		n:     len(runes),
	}
}

func (cit *CharIterator) HasNext() bool {
	return cit.i < cit.n
}

func (cit *CharIterator) Next() (string, error) {

	if !cit.HasNext() {
		return "", errors.New("no more character")
	}

	char := string(cit.runes[cit.i])
	cit.i++

	return char, nil
}

func (cit *CharIterator) JumpTo(tar rune) error {

	if !cit.HasNext() {
		return errors.New("no more character")
	}

	for cit.i = cit.i + 1; cit.HasNext() && cit.runes[cit.i] != tar; cit.i++ {

	}

	if cit.i == cit.n {
		return errors.New("cannot jump to " + strconv.QuoteRune(tar))
	}

	cit.i++

	return nil
}

// i = nextPosition(tar) - 1
func (cit *CharIterator) SkipUntil(tar rune) error {
	err := cit.JumpTo(tar)

	if err != nil {
		return err
	}

	cit.i--
	return nil
}