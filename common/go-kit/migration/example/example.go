package example

import "embed"

//go:embed mysql/*.sql
var Mysql embed.FS
