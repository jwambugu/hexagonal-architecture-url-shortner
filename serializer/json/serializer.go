package json

import (
	"encoding/json"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/shortener"
	errs "github.com/pkg/errors"
)

type Redirect struct {
}

// Decode will unmarshal the bytes to return shortener.Redirect
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}

	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

// Encode returns a bytes slice of shortener.Redirect
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)

	if err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
