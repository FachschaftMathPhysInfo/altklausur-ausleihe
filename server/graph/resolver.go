package graph

import (
	"github.com/adjust/rmq/v3"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB          *gorm.DB
	MinIOClient *minio.Client
	// TODO(chris): should this be a pointer?
	RmqClient rmq.Connection
}
