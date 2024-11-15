package cli

import (
	"buildenv/command"
	"buildenv/config"
	"flag"
	"fmt"
	"os"
	"strings"
)

func newVerifyCmd() *verifyCmd {
	return &verifyCmd{}
}

type verifyCmd struct {
	verify bool
}

func (v *verifyCmd) register() {
	flag.BoolVar(&v.verify, "verify", false, "check and repair toolchain, rootfs, tools and packages for current selected platform.")
}

func (v *verifyCmd) listen() (handled bool) {
	if !v.verify {
		return false
	}

	args := config.VerifyArgs{
		Silent:         silent.silent,
		CheckAndRepair: true,
		BuildType:      buildType.buildType,
	}

	var buildEnvConf config.BuildEnv
	if err := buildEnvConf.Verify(args); err != nil {
		platformName := strings.TrimSuffix(buildEnvConf.Platform, ".json")
		fmt.Print(command.PlatformSelectedFailed(platformName, err))
		os.Exit(1)
	}

	// Silent mode called from buildenv.cmake
	if !silent.silent {
		platformName := strings.TrimSuffix(buildEnvConf.Platform, ".json")
		fmt.Print(command.PlatformSelected(platformName))
	}

	return true
}
