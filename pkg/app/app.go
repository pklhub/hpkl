package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/apple/pkl-go/pkl"
	"hpkl.io/hpkl/pkg/logger"
	"hpkl.io/hpkl/pkg/pklutils"
)

type AppConfig struct {
	Logger          *logger.Logger
	project         *pkl.Project
	ctx             context.Context
	PlainHttp       bool
	CacheDir        string
	DefaultCacheDir string
	WorkingDir      string
	RootDir         string
	Parameters      []string
}

const (
	configPath = ".hpkl/config.pkl"
)

func (a *AppConfig) ProjectOrErr() (*pkl.Project, error) {

	projectFile := filepath.Join(a.WorkingDir, "PklProject")

	if a.project == nil {
		if _, err := os.Stat(projectFile); err == nil {

			proj, err := pklutils.LoadProject(a.ctx, projectFile)

			if err != nil {
				a.Logger.Error("Project file path: %s", projectFile)
				return nil, err
			}
			a.project = proj
		} else {
			return nil, errors.New(fmt.Sprintf("PklProject file not found in the working directory %s", a.WorkingDir))
		}
	}

	return a.project, nil
}

func (a *AppConfig) Project() *pkl.Project {

	p, err := a.ProjectOrErr()
	if err != nil {
		a.Logger.Fatal(err.Error())
	}
	return p
}

func (a *AppConfig) Reset() {
	a.project = nil
}

func NewAppConfig(ctx context.Context, outWriter io.Writer, errWriter io.Writer) (*AppConfig, error) {

	logger := logger.New(outWriter, errWriter)

	return &AppConfig{
		Logger: logger,
		ctx:    ctx,
	}, nil
}
