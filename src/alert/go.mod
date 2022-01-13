module github.com/ca-risken/core/src/alert

go 1.16

require (
	github.com/andybalholm/brotli v1.0.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/ca-risken/common/pkg/database v0.0.0-20211223025030-6bfdc45e906c
	github.com/ca-risken/common/pkg/logging v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/rpc v0.0.0-20220113015330-0e8462d52b5b
	github.com/ca-risken/common/pkg/sqs v0.0.0-20220113015330-0e8462d52b5b // indirect
	github.com/ca-risken/common/pkg/xray v0.0.0-20211118071101-9855266b50a1
	github.com/ca-risken/core/proto/alert v0.0.0-20211207035656-c33310f08a4c
	github.com/ca-risken/core/proto/finding v0.0.0-20211207035656-c33310f08a4c
	github.com/ca-risken/core/proto/project v0.0.0-20211207035656-c33310f08a4c
	github.com/envoyproxy/protoc-gen-validate v0.6.1 // indirect
	github.com/gassara-kys/envconfig v1.4.4
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/jarcoal/httpmock v1.0.6
	github.com/stretchr/testify v1.7.0
	github.com/valyala/fasthttp v1.29.0 // indirect
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210824181836-a4879c3d0e89 // indirect
	google.golang.org/grpc v1.42.0
	gorm.io/driver/mysql v1.1.2 // indirect
	gorm.io/gorm v1.21.13
)
