package repository

import (
	"errors"
	"strings"
)

var ErrorDuplicate = errors.New("dublicate")
var ErrorNotFound = errors.New("not found")

func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
