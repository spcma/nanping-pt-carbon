package id

import (
	"strings"

	"github.com/google/uuid"
)

type Generator interface {
	NumID() int64
}

var Default Generator

func SetDefault(gen Generator) {
	Default = gen
}

func NumID() int64 {
	return Default.NumID()
}

func UUID() string {
	return uuid.New().String()
}

func SimpleUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
