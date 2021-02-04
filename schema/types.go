package schema

import "github.com/yaoapp/xun/dbal/schema"

// Blueprint a blueprint of how the database is constructed
type Blueprint struct{ schema.Blueprint }

// Schema The database Schema interface
type Schema interface{ schema.Schema }
