package cli

import (
	"buildenv/config"
	"buildenv/console"
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
	flag.BoolVar(&v.verify, "verify", false, "verify buildenv")
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
		fmt.Print(console.PlatformSelectedFailed(platformName, err))
		os.Exit(1)
	}

	// Silent mode called from buildenv.cmake
	if !silent.silent {
		platformName := strings.TrimSuffix(buildEnvConf.Platform, ".json")
		fmt.Print(console.PlatformSelected(platformName))
	}

	return true
}
