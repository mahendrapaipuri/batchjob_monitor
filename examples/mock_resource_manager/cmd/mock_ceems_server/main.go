//go:build cgo
// +build cgo

// Boiler plate code to create a new instance of CEEMSServer entrypoint
package main

// Ensure to import the current mock scheduler
import (
	"fmt"
	"os"

	_ "github.com/mahendrapaipuri/ceems/examples/mock_resource_manager/pkg/resource"

	"github.com/mahendrapaipuri/ceems/pkg/api/cli"

	// If existing schedulers in CEEMS are needed, they need to be imported too
	// For instance to import slurm manager, following import statement must be added
	_ "github.com/mahendrapaipuri/ceems/pkg/api/resource/slurm"
)

// Main entry point for `usagestats` app
func main() {
	// Create a new app
	usageStatsServerApp, err := cli.NewCEEMSServer()
	if err != nil {
		panic("Failed to create an instance of Usage Stats Server App")
	}

	// Main entrypoint of the app
	if err := usageStatsServerApp.Main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
