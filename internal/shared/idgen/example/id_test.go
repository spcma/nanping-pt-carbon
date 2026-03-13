package example

import (
	"app/internal/common/id"
	"log"
	"testing"
)

func Test1(t *testing.T) {
	workId := 1
	generateFactory, err := id.NewIdgenGenerateFactory(int64(workId))
	if err != nil {
		log.Fatal(err)
	}
	id.SetDefault(generateFactory)

	log.Println(id.NumID())
	log.Println(id.NumID())
	log.Println(id.NumID())
	log.Println(id.NumID())
	log.Println(id.NumID())

	log.Println(id.SimpleUUID())
	log.Println(id.UUID())
}
