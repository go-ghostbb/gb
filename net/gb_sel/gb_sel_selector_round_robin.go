package gbsel

import (
	"context"
	"github.com/Ghostbb-io/gb/internal/intlog"
	"sync"
)

type selectorRoundRobin struct {
	mu    sync.Mutex
	nodes Nodes
	next  int
}

func NewSelectorRoundRobin() Selector {
	return &selectorRoundRobin{
		nodes: make(Nodes, 0),
	}
}

func (s *selectorRoundRobin) Update(ctx context.Context, nodes Nodes) error {
	intlog.Printf(ctx, `Update nodes: %s`, nodes.String())
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = nodes
	s.next = 0
	return nil
}

func (s *selectorRoundRobin) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.nodes) == 0 {
		return
	}
	node = s.nodes[s.next]
	s.next = (s.next + 1) % len(s.nodes)
	intlog.Printf(ctx, `Picked node: %s`, node.Address())
	return
}
