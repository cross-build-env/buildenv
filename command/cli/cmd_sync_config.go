package cli

import (
	"buildenv/command"
	"buildenv/config"
	"buildenv/pkg/io"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func newSyncConfigCmd() *syncConfigCmd {
	return &syncConfigCmd{}
}

type syncConfigCmd struct {
	sync bool
}

func (s *syncConfigCmd) register() {
	flag.BoolVar(&s.sync, "sync", false, "create buildenv.json or sync conf repo defined in buildenv.json.")
}

func (s *syncConfigCmd) listen() (handled bool) {
	if !s.sync {
		return false
	}

	// Create buildenv.json if not exist.
	confPath := filepath.Join(config.Dirs.WorkspaceDir, "buildenv.json")
	if !io.PathExists(confPath) {
		if err := os.MkdirAll(filepath.Dir(confPath), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		var buildenv config.BuildEnv
		buildenv.JobNum = runtime.NumCPU()

		bytes, err := json.MarshalIndent(buildenv, "", "    ")
		if err != nil {
			fmt.Print(command.SyncFailed(err))
			os.Exit(1)
		}
		if err := os.WriteFile(confPath, []byte(bytes), os.ModePerm); err != nil {
			fmt.Print(command.SyncFailed(err))
			os.Exit(1)
		}

		fmt.Print(command.SyncSuccess(false))
		return false
	}

	// Sync conf repo with repo url.
	bytes, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Print(command.SyncFailed(err))
		os.Exit(1)
	}

	// Unmarshall with buildenv.json.
	var buildenv config.BuildEnv
	if err := json.Unmarshal(bytes, &buildenv); err != nil {
		fmt.Print(command.SyncFailed(err))
		os.Exit(1)
	}

	// Sync repo.
	output, err := buildenv.SyncRepo(buildenv.ConfRepo, buildenv.ConfRepoRef)
	if err != nil {
		fmt.Print(command.SyncFailed(err))
		os.Exit(1)
	}

	fmt.Println(output)
	fmt.Print(command.SyncSuccess(true))

	return true
}
