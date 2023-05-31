package localdbconnector

import (
	"github.com/tzapio/tzap/pkg/embed/localdb"
	"github.com/tzapio/tzap/pkg/types"
)

type LocalembedTGenerator struct {
	*types.UnimplementedTGenerator
	db *localdb.FileDB[types.Vector]
}

func InitiateLocalDB(filePath string) (types.TGenerator, error) {
	db, err := localdb.NewFileDB[types.Vector](filePath)
	if err != nil {
		return nil, err
	}
	return &LocalembedTGenerator{db: db}, nil
}
