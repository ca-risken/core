module github.com/ca-risken/core/src/project

go 1.16

require (
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/database v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/rpc v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/pkg/model v0.0.0-20210906100342-c1bbb08cc3e4
	github.com/ca-risken/core/proto/iam v0.0.0-20210906100342-c1bbb08cc3e4
	github.com/ca-risken/core/proto/project v0.0.0-20210906100342-c1bbb08cc3e4
	github.com/envoyproxy/protoc-gen-validate v0.6.1 // indirect
	github.com/gassara-kys/envconfig v1.4.4
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210824181836-a4879c3d0e89 // indirect
	google.golang.org/grpc v1.42.0
	gorm.io/gorm v1.21.12
)
