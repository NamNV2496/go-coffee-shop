package configs

import (
	"github.com/google/wire"
)

var ConfigWireSet = wire.NewSet(
	GetConfigFromYaml,
	wire.FieldsOf(new(Config), "Grpc"),
	wire.FieldsOf(new(Config), "Database"),
	wire.FieldsOf(new(Config), "Redis"),
	wire.FieldsOf(new(Config), "Kafka"),
	wire.FieldsOf(new(Config), "Cron"),
	wire.FieldsOf(new(Config), "S3"),
)
