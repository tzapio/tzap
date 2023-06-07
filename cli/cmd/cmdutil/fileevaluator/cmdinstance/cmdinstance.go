package cmdinstance

import (
	"os"

	"github.com/tzapio/tzap/pkg/project"
)

func NewProjectDB() (project.ProjectDB2, error) {
	projectDB := project.ProjectDB2{}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	localProject, err := NewLocalProject(cwd)
	if err != nil {
		return nil, err
	}
	projectDB[project.LOCALPROJECTNAME] = localProject
	return projectDB, nil
}
