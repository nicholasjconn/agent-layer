package clients

import (
	"os"

	"github.com/nicholasjconn/agent-layer/internal/config"
	"github.com/nicholasjconn/agent-layer/internal/run"
	"github.com/nicholasjconn/agent-layer/internal/sync"
)

// LaunchFunc launches a client after sync and run setup.
type LaunchFunc func(project *config.ProjectConfig, runInfo *run.Info, env []string) error

// EnabledSelector returns the enabled flag for a client.
type EnabledSelector func(cfg *config.Config) *bool

// Run performs the standard client launch pipeline: load config, sync, create run dir, launch.
func Run(root string, name string, enabled EnabledSelector, launch LaunchFunc) error {
	project, err := config.LoadProjectConfig(root)
	if err != nil {
		return err
	}
	if err := sync.EnsureEnabled(name, enabled(&project.Config)); err != nil {
		return err
	}

	if err := sync.RunWithProject(root, project); err != nil {
		return err
	}

	runInfo, err := run.Create(root)
	if err != nil {
		return err
	}

	env := BuildEnv(os.Environ(), project.Env, runInfo)

	return launch(project, runInfo, env)
}
