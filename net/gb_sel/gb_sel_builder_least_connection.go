package gbsel

type builderLeastConnection struct{}

func NewBuilderLeastConnection() Builder {
	return &builderLeastConnection{}
}

func (*builderLeastConnection) Name() string {
	return "BalancerLeastConnection"
}

func (*builderLeastConnection) Build() Selector {
	return NewSelectorLeastConnection()
}
