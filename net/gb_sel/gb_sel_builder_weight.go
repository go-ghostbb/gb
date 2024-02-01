package gbsel

type builderWeight struct{}

func NewBuilderWeight() Builder {
	return &builderWeight{}
}

func (*builderWeight) Name() string {
	return "BalancerWeight"
}

func (*builderWeight) Build() Selector {
	return NewSelectorWeight()
}
