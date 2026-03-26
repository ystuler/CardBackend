package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON(schema interface{}, raw []byte) error {
	if len(raw) == 0 {
		return io.ErrUnexpectedEOF
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(schema); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return &json.SyntaxError{Offset: 0}
	}

	return nil
}

func DecodeJSONRequest(r *http.Request, schema interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := DecodeJSON(schema, body); err != nil {
		return err
	}

	return nil
}
