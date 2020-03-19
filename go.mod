module github.com/shine-o/shine.engine.world

go 1.13

require (
	github.com/go-redis/redis/v7 v7.2.0
	github.com/google/logger v1.0.1
	github.com/google/uuid v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/shine-o/shine.engine.networking v0.0.0-20200319001014-56ac0080099d
	github.com/shine-o/shine.engine.protocol-buffers/login-world v0.0.0-20200314174953-2f346ed277bf
	github.com/shine-o/shine.engine.structs v0.0.0-20200317164417-264763e42420
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	google.golang.org/genproto v0.0.0-20200318110522-7735f76e9fa5 // indirect
	google.golang.org/grpc v1.28.0
)
