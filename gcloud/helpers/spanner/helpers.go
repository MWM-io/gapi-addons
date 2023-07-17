package spanner

import (
	"fmt"

	"cloud.google.com/go/spanner"
)

// GetDBPath is a simple function that returns the database path to use as
// argument for spanner.NewClients
func GetDBPath(projectID, instanceName, dbName string) string {
	return fmt.Sprintf("projects/%s/instances/%s/databases/%s",
		projectID,
		instanceName,
		dbName)
}

// ToSpannerKey wrap given arguments and return a spanner.Key
func ToSpannerKey(args ...interface{}) spanner.Key {
	return args
}
