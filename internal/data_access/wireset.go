package database

import (
	"github.com/google/wire"
)

var DataWireSet = wire.NewSet(
	// NewDatabaseConfig,
	InitializeAndMigrateUpDB,
	InitializeGoquDB,
	NewMigrator,
)
