module github.com/tzapio/tzap/pkg/connectors/googlevoiceconnector

go 1.20

replace github.com/tzapio/tzap => ../../../

replace github.com/tzapio/tzap/pkg/tzapconnect => ../../tzapconnect

require (
	cloud.google.com/go/speech v1.15.0
	cloud.google.com/go/texttospeech v1.6.0
	github.com/tzapio/tzap v0.7.11
	github.com/tzapio/tzap/pkg/tzapconnect v0.7.11
	google.golang.org/api v0.120.0
)

require (
	cloud.google.com/go v0.110.0 // indirect
	cloud.google.com/go/compute v1.19.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/longrunning v0.4.1 // indirect
	github.com/dlclark/regexp2 v1.9.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/s2a-go v0.1.3 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.8.0 // indirect
	github.com/sashabaranov/go-openai v1.9.0 // indirect
	github.com/tiktoken-go/tokenizer v0.1.0 // indirect
	github.com/tzapio/tzap/pkg/connectors/openaiconnector v0.7.11 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/oauth2 v0.7.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.54.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
