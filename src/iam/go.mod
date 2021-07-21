module github.com/CyberAgent/mimosa-core/src/iam

go 1.16

require (
	github.com/CyberAgent/mimosa-common/pkg/database v0.0.0-20210721063343-44cefe7f590e
	github.com/CyberAgent/mimosa-common/pkg/xray v0.0.0-20210721063343-44cefe7f590e
	github.com/CyberAgent/mimosa-core/pkg/model v0.0.0-20210712023706-882d5424f2f1
	github.com/CyberAgent/mimosa-core/proto/iam v0.0.0-20210712023706-882d5424f2f1
	github.com/aws/aws-xray-sdk-go v1.6.0
	github.com/go-ozzo/ozzo-validation/v4 v4.3.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210708141623-e76da96a951f // indirect
	google.golang.org/grpc v1.39.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.12
)
