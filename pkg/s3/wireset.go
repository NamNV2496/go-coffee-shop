package s3

import (
	"github.com/google/wire"
)

var FileWireSet = wire.NewSet(
	NewS3Client,
)
