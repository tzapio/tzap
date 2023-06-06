package localdbconnector

import (
	"path"

	"github.com/tzapio/tzap/internal/logging/tl"
	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/project"
	"github.com/tzapio/tzap/pkg/types"
)

type LocalembedTGenerator struct {
	*types.UnimplementedTGenerator
	dbs map[project.ProjectName]*localdb.FileDB[types.Vector]
}

func InitiateLocalDB(projectDB project.ProjectDB) (types.TGenerator, error) {
	dbs := make(map[project.ProjectName]*localdb.FileDB[types.Vector])
	for projectName, projectDir := range projectDB {
		tl.Logger.Println("Initiating fileembeddings localDB for project", projectName)
		db, err := localdb.NewFileDB[types.Vector](path.Join(string(projectDir), "fileembeddings.db"))
		if err != nil {
			return nil, err
		}
		dbs[projectName] = db
	}
	return &LocalembedTGenerator{dbs: dbs}, nil
}
