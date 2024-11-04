package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type BuildEnvConf struct {
	ResRepoUrl  string `json:"res_repo_url"`
	ConfRepoUrl string `json:"conf_repo_url"`
	Platform    string `json:"platform"`
}

func (b *BuildEnvConf) Verify(onlyFields bool) error {
	bytes, err := os.ReadFile("conf/buildenv.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, b); err != nil {
		return err
	}

	if b.Platform == "" {
		return fmt.Errorf("platform is empty")
	}

	filePath := filepath.Join(PlatformsDir, b.Platform)
	if !pathExists(filePath) {
		return fmt.Errorf("platform file not exists: %s", filePath)
	}

	var buildenv BuildEnv
	if err := buildenv.Read(filePath); err != nil {
		return err
	}

	if err := buildenv.Verify(b.ResRepoUrl, onlyFields); err != nil {
		return err
	}

	return nil
}
