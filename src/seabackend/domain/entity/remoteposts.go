package entity

import (
	"encoding/json"
	"strings"
)

type Field int

const (
	FieldTitle Field = iota
	FieldBody
	FieldAll
)

type RemotePost struct {
	UserID json.Number `json:"userId"`
	ID     json.Number `json:"id"`
	Title  string      `json:"title,omitempty"`
	Body   string      `json:"body,omitempty"`
}

func (rp RemotePost) Contains(query string, field Field) bool {
	if query == "" {
		return true
	}

	switch field {
	case FieldTitle:
		if strings.Contains(strings.ToLower(rp.Title), strings.ToLower(query)) {
			return true
		}
	case FieldBody:
		if strings.Contains(strings.ToLower(rp.Body), strings.ToLower(query)) {
			return true
		}
	case FieldAll:
		if strings.Contains(strings.ToLower(rp.Title), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(rp.Body), strings.ToLower(query)) {
			return true
		}
	}

	return false
}
