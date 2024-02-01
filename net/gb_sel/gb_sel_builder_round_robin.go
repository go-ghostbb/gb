package gbsel

type builderRoundRobin struct{}

func NewBuilderRoundRobin() Builder {
	return &builderRoundRobin{}
}

func (*builderRoundRobin) Name() string {
	return "BalancerRoundRobin"
}

func (*builderRoundRobin) Build() Selector {
	return NewSelectorRoundRobin()
}
