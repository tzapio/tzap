package cmdinstance

import (
	"path"

	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator"
	"github.com/tzapio/tzap/cli/cmd/cmdutil/fileevaluator/localwalker"
	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

type LocalProject struct {
	projectName              project.ProjectName
	projectDir               string
	baseDir                  string
	embeddingCollection      types.DBCollectionInterface[types.Vector]
	filestampsDB             types.DBCollectionInterface[int64]
	embeddingCacheDB         types.DBCollectionInterface[string]
	*localwalker.LocalWalker //GetFiles() @TODO: Refactor to FS interface?
}

func NewFilestampCache(projectDir project.ProjectDir) (types.DBCollectionInterface[int64], error) {
	return localdb.NewFileDB[int64](path.Join(string(projectDir), "filesTimestamps.db"))
}
func NewEmbeddingsCache(projectDir project.ProjectDir) (types.DBCollectionInterface[string], error) {
	return localdb.NewFileDB[string](path.Join(string(projectDir), "embeddingsCache.db"))
}
func NewEmbeddingsCollection(projectDir project.ProjectDir) (types.DBCollectionInterface[types.Vector], error) {
	return localdb.NewFileDB[types.Vector](path.Join(string(projectDir), "fileembeddings.db"))
}
func NewLocalProject(baseDir string) (project.Project, error) {
	filesStampsDB, err := NewFilestampCache("./.tzap-data")
	if err != nil {
		panic(err)
	}

	embeddingCacheDB, err := NewEmbeddingsCache("./.tzap-data")
	if err != nil {
		panic(err)
	}
	projectDir := path.Join(baseDir, "./.tzap-data")
	fileevaluator, err := fileevaluator.New(baseDir)
	if err != nil {
		return nil, err
	}
	localWalker := localwalker.New(fileevaluator, baseDir, baseDir)
	embeddingCollection, err := localdb.NewFileDB[types.Vector](path.Join(projectDir, "fileembeddings.db"))
	if err != nil {
		return nil, err
	}
	localProject := &LocalProject{
		projectName:         project.LOCALPROJECTNAME,
		baseDir:             baseDir,
		filestampsDB:        filesStampsDB,
		embeddingCacheDB:    embeddingCacheDB,
		projectDir:          projectDir,
		embeddingCollection: embeddingCollection,
		LocalWalker:         localWalker,
	}

	return localProject, nil
}
func (l *LocalProject) CanIndex() bool {
	return true
}
func (l *LocalProject) GetTimestampCache() types.DBCollectionInterface[int64] {
	return l.filestampsDB
}

func (l *LocalProject) GetEmbeddingsCache() types.DBCollectionInterface[string] {
	return l.embeddingCacheDB
}

// GetEmbedding implements project.Project
func (l *LocalProject) GetEmbeddingCollection() types.DBCollectionInterface[types.Vector] {
	return l.embeddingCollection
}

// GetProjectName implements project.Project
func (l *LocalProject) GetProjectName() project.ProjectName {
	return l.projectName
}
