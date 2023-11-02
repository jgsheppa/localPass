package models

import "strings"

const (
	ErrNotFound   modelError   = "models: resource not found"
	ErrIDInvalid  privateError = "models: ID provided was invalid"
	ErrURLInvalid modelError   = "Invalid URL provided"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
