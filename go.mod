module github.com/shine-o/shine.engine.zone

go 1.13

require (
	github.com/RoaringBitmap/roaring v0.4.23
	github.com/go-pg/pg/v9 v9.1.6
	github.com/go-redis/redis/v7 v7.2.0
	github.com/google/logger v1.1.0
	github.com/google/uuid v1.1.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/shine-o/shine.engine.core v0.0.3
	github.com/shine-o/shine.engine.world v0.0.0-00010101000000-000000000000
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	go.etcd.io/bbolt v1.3.4
	golang.org/x/image v0.0.0-20191009234506-e7c1f5e7dbb8
	google.golang.org/grpc v1.29.0
)

replace github.com/shine-o/shine.engine.core => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.core

replace github.com/shine-o/shine.engine.world => C:\Users\marbo\go\src\github.com\shine-o\shine.engine.world
