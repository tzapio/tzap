module github.com/tzapio/tzap/example

go 1.20

replace github.com/tzapio/tzap => ../

replace github.com/tzapio/tzap/pkg/connectors/openaiconnector => ../pkg/connectors/openaiconnector

replace github.com/tzapio/tzap/pkg/connectors/redisembeddbconnector => ../pkg/connectors/redisembeddbconnector

replace github.com/tzapio/tzap/pkg/tzapconnect => ../pkg/tzapconnect

replace github.com/tzapio/tzap/pkg/tzapaction => ../pkg/tzapaction

replace github.com/tzapio/tzap/cli => ../cli
