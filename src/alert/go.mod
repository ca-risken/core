module github.com/CyberAgent/mimosa-core/src/alert

go 1.15

require (
	github.com/CyberAgent/mimosa-core/pkg/model v0.0.0-20201016035408-d7ace0dc922a
	github.com/CyberAgent/mimosa-core/proto/alert v0.0.0-20201016035408-d7ace0dc922a
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/jarcoal/httpmock v1.0.6
	github.com/jinzhu/gorm v1.9.16
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/vikyd/zero v0.0.0-20190921142904-0f738d0bc858
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb // indirect
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect
	google.golang.org/genproto v0.0.0-20201015140912-32ed001d685c // indirect
	google.golang.org/grpc v1.33.0
)
