package db

import "os"

var (
	datastoreName = os.Getenv("POSTGRES_CONNECTION")
)
