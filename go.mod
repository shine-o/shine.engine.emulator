module github.com/shine-o/shine.engine.login

go 1.13

require (
	github.com/go-redis/redis/v7 v7.2.0
	github.com/google/logger v1.0.1
	github.com/google/uuid v1.1.1
	github.com/jinzhu/gorm v1.9.12
	github.com/lib/pq v1.3.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shine-o/shine.engine.networking v0.0.0-20200317193756-3f4ef9b27934
	github.com/shine-o/shine.engine.protocol-buffers/login-world v0.0.0-20200314174953-2f346ed277bf
	github.com/shine-o/shine.engine.structs v0.0.0-20200317164417-264763e42420
	github.com/shine-o/shine.engine.world v0.0.0-20200317162035-b5df72e4dc88
	github.com/shine-o/shine.engine.world/service v0.0.0-20200314172549-7283ba661f79 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	google.golang.org/genproto v0.0.0-20200317114155-1f3552e48f24 // indirect
	google.golang.org/grpc v1.28.0
)

replace github.com/shine-o/shine.engine.networking => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.networking
