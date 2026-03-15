package id

import (
	"sync"

	"github.com/yitter/idgenerator-go/idgen"
)

var idgenOnce sync.Once

type IdgenGenerateFactory struct {
}

func NewIdgenGenerateFactory(workId int64) (*IdgenGenerateFactory, error) {
	idgenOnce.Do(func() {
		idgen.SetIdGenerator(idgen.NewIdGeneratorOptions(uint16(workId)))
	})

	return &IdgenGenerateFactory{}, nil
}

func (g *IdgenGenerateFactory) NumID() int64 {
	return idgen.NextId()
}
