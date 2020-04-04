module github.com/shine-o/shine.engine.login

go 1.13

require (
	github.com/go-pg/pg/v9 v9.1.5
	github.com/go-redis/redis/v7 v7.2.0
	github.com/google/logger v1.0.1
	github.com/google/uuid v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/shine-o/shine.engine.networking v0.0.0-20200401184904-1c8aadb06909
	github.com/shine-o/shine.engine.protocol-buffers/login-world v0.0.0-20200329210513-1b648aeb6624
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	golang.org/x/net v0.0.0-20200320220750-118fecf932d8 // indirect
	google.golang.org/genproto v0.0.0-20200319113533-08878b785e9c // indirect
	google.golang.org/grpc v1.28.0
)

replace github.com/shine-o/shine.engine.networking => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.networking
