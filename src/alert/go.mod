module github.com/CyberAgent/mimosa-core/src/alert

go 1.16

require (
	github.com/CyberAgent/mimosa-common/pkg/xray v0.0.0-20210709182517-c12a4e8eed4d
	github.com/CyberAgent/mimosa-core/pkg/model v0.0.0-20210628032046-7e6a43522da4
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20210628032046-7e6a43522da4
	github.com/CyberAgent/mimosa-core/proto/finding v0.0.0-20210628032046-7e6a43522da4
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/jarcoal/httpmock v1.0.6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.6.1
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
	google.golang.org/genproto v0.0.0-20210624195500-8bfb893ecb84 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.0 // indirect
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)
