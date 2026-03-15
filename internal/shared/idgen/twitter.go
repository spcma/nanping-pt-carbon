package id

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var twitterOnce sync.Once

type TwitterGenerateFactory struct {
	workId int64
	node   *snowflake.Node
}

func NewTwitterGenerateFactory(workId int64) (*TwitterGenerateFactory, error) {
	var node *snowflake.Node
	var err error

	twitterOnce.Do(func() {
		node, err = snowflake.NewNode(workId)
	})
	if err != nil {
		return nil, err
	}

	return &TwitterGenerateFactory{
		node: node,
	}, nil
}

func (g *TwitterGenerateFactory) NumID() int64 {
	return g.node.Generate().Int64()
}
