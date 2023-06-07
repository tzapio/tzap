package cmdinstance

import (
	"path"

	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator/zipwalker"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

type ZipProject struct {
	projectName          project.ProjectName
	projectDir           project.ProjectDir
	baseDir              string
	embeddingCollection  types.DBCollectionInterface[types.Vector]
	*zipwalker.ZipWalker //GetFiles() @TODO: Refactor to FS interface?
}

func NewLocalZipProject(name, relativeDirInZip, url string, excludePatterns, includePatterns []string) (project.Project, error) {
	fileevaluator := fileevaluator.NewWithPatterns(excludePatterns, includePatterns)
	projectDir := project.ProjectDir(path.Join("./.tzap-data", name))
	zipwalker := zipwalker.New(fileevaluator, relativeDirInZip, url)

	embeddingCollection, err := NewEmbeddingsDB(projectDir)
	if err != nil {
		return nil, err
	}
	zipProject := &ZipProject{
		projectName:         project.LOCALPROJECTNAME,
		baseDir:             relativeDirInZip,
		projectDir:          projectDir,
		embeddingCollection: embeddingCollection,
		ZipWalker:           zipwalker,
	}
	return zipProject, nil
}
func (l *ZipProject) CanIndex() bool {
	return true
}

// GetEmbeddingsCache implements project.Project
func (*ZipProject) GetEmbeddingsCache() types.DBCollectionInterface[string] {
	panic("unimplemented")
}

// GetTimestampCache implements project.Project
func (*ZipProject) GetTimestampCache() types.DBCollectionInterface[int64] {
	panic("unimplemented")
}

// GetEmbedding implements project.Project
func (l *ZipProject) GetEmbeddingCollection() types.DBCollectionInterface[types.Vector] {
	return l.embeddingCollection
}

// GetProjectName implements project.Project
func (l *ZipProject) GetProjectName() project.ProjectName {
	return l.projectName
}
