package cmdinstance

import (
	"os"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/project"
)

func LoadProject(baseDir string, lib string) (project.Project, error) {
	var projectP project.Project
	if lib != "" {
		var name project.ProjectName = project.ProjectName(lib)
		libProject, err := NewLocalLibProject(baseDir, name)
		if err != nil {
			return nil, err
		}
		tl.Logger.Println("Loaded lib ProjectDB:", name, libProject)
		projectP = libProject
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		localProject, err := NewLocalProject(cwd)
		if err != nil {
			return nil, err
		}
		projectP = localProject
	}

	return projectP, nil
}
