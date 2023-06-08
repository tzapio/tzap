package project

import (
	"context"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/types"
)

type ProjectDir string
type ProjectName string

type ProjectDB map[ProjectName]ProjectDir
type ProjectDB2 map[ProjectName]Project

const LOCALPROJECTNAME ProjectName = "@LOCAL"

type LocalProject struct {
	ProjectName ProjectName
	ProjectDir  ProjectDir
}

type Project interface {
	GetProjectName() ProjectName
	GetFiles() ([]types.FileReader, error)
	GetEmbeddingCollection() types.DBCollectionInterface[types.Vector]
	GetTimestampCache() types.DBCollectionInterface[int64]
	GetEmbeddingsCache() types.DBCollectionInterface[string]
	CanIndex() bool
}
type EmbedStore interface {
}

var projectKey = struct{ projectKey string }{}

func SetProjectInContext(ctx context.Context, project Project) context.Context {
	tl.DeepLogger.Println("SetProjectInContext:", project.GetProjectName())
	return context.WithValue(ctx, projectKey, project)
}
func GetProjectFromContext(ctx context.Context) Project {
	project := ctx.Value(projectKey).(Project)
	tl.DeepLogger.Println("GetProjectFromContext:", project.GetProjectName())
	return project
}
