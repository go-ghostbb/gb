package gbsel

// defaultBuilder is the default Builder for globally used purpose.
var defaultBuilder = NewBuilderRoundRobin()

// SetBuilder sets the default builder for globally used purpose.
func SetBuilder(builder Builder) {
	defaultBuilder = builder
}

// GetBuilder returns the default builder for globally used purpose.
func GetBuilder() Builder {
	return defaultBuilder
}
