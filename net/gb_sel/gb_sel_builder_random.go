package gbsel

type builderRandom struct{}

func NewBuilderRandom() Builder {
	return &builderRandom{}
}

func (*builderRandom) Name() string {
	return "BalancerRandom"
}

func (*builderRandom) Build() Selector {
	return NewSelectorRandom()
}
