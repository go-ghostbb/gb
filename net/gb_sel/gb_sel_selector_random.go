package gbsel

import (
	"context"
	"ghostbb.io/internal/intlog"
	gbrand "ghostbb.io/util/gb_rand"
	"sync"
)

type selectorRandom struct {
	mu    sync.RWMutex
	nodes Nodes
}

func NewSelectorRandom() Selector {
	return &selectorRandom{
		nodes: make([]Node, 0),
	}
}

func (s *selectorRandom) Update(ctx context.Context, nodes Nodes) error {
	intlog.Printf(ctx, `Update nodes: %s`, nodes.String())
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = nodes
	return nil
}

func (s *selectorRandom) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.nodes) == 0 {
		return nil, nil, nil
	}
	node = s.nodes[gbrand.Intn(len(s.nodes))]
	intlog.Printf(ctx, `Picked node: %s`, node.Address())
	return node, nil, nil
}
