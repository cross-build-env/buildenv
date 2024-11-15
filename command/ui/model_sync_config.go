package ui

import (
	"buildenv/command"
	"buildenv/config"
	"buildenv/pkg/color"
	"buildenv/pkg/io"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
)

func newSyncConfigModel(goback func()) *syncConfigModel {
	content := fmt.Sprintf("\nClone or synch repo of conf.\n"+
		"-----------------------------------\n"+
		"%s.\n\n"+
		"%s",
		color.Sprintf(color.Blue, "This will create a buildenv.json if not exist, otherwise it'll checkout the latest conf repo with specified repo REF"),
		color.Sprintf(color.Gray, "[↵ -> execute | ctrl+c/q -> quit]"))

	return &syncConfigModel{
		content: content,
		goback:  goback,
	}
}

type syncConfigModel struct {
	content string
	goback  func()
}

func (s syncConfigModel) Init() tea.Cmd {
	return nil
}

func (s syncConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit

		case "enter":
			if output, err := s.syncRepo(); err != nil {
				s.content += "\r" + color.Sprintf(color.Red, err.Error())
			} else {
				s.content += "\r" + output + "\n" + command.SyncSuccess(true)
			}
			return s, tea.Quit

		case "esc":
			s.goback()
			return s, nil
		}
	}
	return s, nil
}

func (s syncConfigModel) View() string {
	return s.content
}

func (s syncConfigModel) syncRepo() (string, error) {
	// Create buildenv.json if not exist.
	confPath := filepath.Join(config.Dirs.WorkspaceDir, "buildenv.json")
	if !io.PathExists(confPath) {
		if err := os.MkdirAll(filepath.Dir(confPath), os.ModePerm); err != nil {
			return "", err
		}

		var buildenv config.BuildEnv
		buildenv.JobNum = runtime.NumCPU()

		bytes, err := json.MarshalIndent(buildenv, "", "    ")
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(confPath, []byte(bytes), os.ModePerm); err != nil {
			return "", err
		}

		return command.SyncSuccess(false), nil
	}

	// Sync conf repo with repo url.
	bytes, err := os.ReadFile(confPath)
	if err != nil {
		return "", err
	}

	// Unmarshall with buildenv.json.
	var buildenv config.BuildEnv
	if err := json.Unmarshal(bytes, &buildenv); err != nil {
		return "", err
	}

	// Sync repo.
	return buildenv.SyncRepo(buildenv.ConfRepo, buildenv.ConfRepoRef)
}
