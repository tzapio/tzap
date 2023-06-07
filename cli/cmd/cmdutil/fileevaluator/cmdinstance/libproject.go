package cmdinstance

import (
	"path"

	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

type LibProject struct {
	projectName         project.ProjectName
	projectDir          string
	baseDir             string
	embeddingCollection types.DBCollectionInterface[types.Vector]
}

func NewLocalLibProject(baseDir string, name project.ProjectName) (project.Project, error) {
	projectDir := path.Join(baseDir, "./.tzap-data", string(name))

	embeddingCollection, err := localdb.NewFileDB[types.Vector](path.Join(projectDir, "fileembeddings.db"))
	if err != nil {
		return nil, err
	}
	localProject := &LibProject{
		projectName:         name,
		baseDir:             baseDir,
		projectDir:          projectDir,
		embeddingCollection: embeddingCollection,
	}
	return localProject, nil
}
func (l *LibProject) CanIndex() bool {
	return false
}

// GetEmbeddingsCache implements project.Project
func (*LibProject) GetEmbeddingsCache() types.DBCollectionInterface[string] {
	panic("Local LibProject does not implement GetEmbeddingsCache() - Do not index libproject")
}

// GetTimestampCache implements project.Project
func (*LibProject) GetTimestampCache() types.DBCollectionInterface[int64] {
	panic("Local LibProject does not implement GetTimestampCache() - Do not index libproject")
}

func (l *LibProject) GetEmbeddingCollection() types.DBCollectionInterface[types.Vector] {
	return l.embeddingCollection
}

// GetFiles implements project.Project
func (*LibProject) GetFiles() ([]types.FileReader, error) {
	panic("Local LibProject does not implement GetFiles() - Do not index libproject")
}

// GetProjectName implements project.Project
func (l *LibProject) GetProjectName() project.ProjectName {
	return l.projectName
}
