package cmdinstance

import (
	"path"

	"github.com/tzapio/tzap/cli/cmd/cmdinstance/zipwalker"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

type ZipProject struct {
	projectName          project.ProjectName
	projectDir           project.ProjectDir
	baseDir              string
	embeddingCollection  types.DBCollectionInterface[types.Vector]
	embeddingsCache      types.DBCollectionInterface[string]
	filestampsCache      types.DBCollectionInterface[int64]
	*zipwalker.ZipWalker //GetFiles() @TODO: Refactor to FS interface?
}

func NewLocalZipProject(name, relativeDirInZip, url string) (project.Project, error) {
	tl.Logger.Println("URL: ", url)
	fileevaluator := fileevaluator.NewWithBasePatterns()
	projectDir := project.ProjectDir(path.Join("./.tzap-data", name))
	zipwalker := zipwalker.New(fileevaluator, relativeDirInZip, url)

	embeddingCollection, err := NewEmbeddingsCollection(projectDir)
	if err != nil {
		return nil, err
	}
	embeddingsCache, err := NewEmbeddingsCache(projectDir)
	if err != nil {
		return nil, err
	}
	filestampCache, err := NewFilestampCache(projectDir)
	if err != nil {
		return nil, err
	}
	zipProject := &ZipProject{
		projectName:         project.LOCALPROJECTNAME,
		baseDir:             relativeDirInZip,
		projectDir:          projectDir,
		embeddingCollection: embeddingCollection,
		embeddingsCache:     embeddingsCache,
		filestampsCache:     filestampCache,
		ZipWalker:           zipwalker,
	}
	return zipProject, nil
}
func (l *ZipProject) CanIndex() bool {
	return true
}

// GetEmbeddingsCache implements project.Project
func (l *ZipProject) GetEmbeddingsCache() types.DBCollectionInterface[string] {
	return l.embeddingsCache
}

// GetTimestampCache implements project.Project
func (l *ZipProject) GetTimestampCache() types.DBCollectionInterface[int64] {
	return l.filestampsCache
}

// GetEmbedding implements project.Project
func (l *ZipProject) GetEmbeddingCollection() types.DBCollectionInterface[types.Vector] {
	return l.embeddingCollection
}

// GetProjectName implements project.Project
func (l *ZipProject) GetProjectName() project.ProjectName {
	return l.projectName
}
