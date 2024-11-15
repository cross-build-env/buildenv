package cli

import "flag"

func newBuildTypeCmd() *buildTypeCmd {
	return &buildTypeCmd{}
}

type buildTypeCmd struct {
	buildType string
}

func (b *buildTypeCmd) register() {
	flag.StringVar(&b.buildType, "build_type", "Release", "value of CMAKE_BUILD_TYPE")
}