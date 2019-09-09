package util

import (
	"encoding/json"
	"net/url"
)

const (
	InvalidUrl       = "not a valid url"
	NotFoundUrl      = "shortened url not found"
	AlreadyExistsUrl = "custom alias already exists"
)

type Error struct {
	Alias       string `json:"alias"`
	ErrCode     string `json:"err_code"`
	Description string `json:"description"`
}

func (c Error) Error() string {
	bytes, _ := json.Marshal(&c)
	return string(bytes)
}

func NewError(alias, errorCode, description string) error {
	return Error{
		Alias:       alias,
		ErrCode:     errorCode,
		Description: description,
	}
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
