//go:build mage
// +build mage

package main

import (
	"os"
	"path/filepath"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/magefile/mage/mg"
)

func init() {
	// Set go version for docker builds
	os.Setenv("GO_VERSION", "1.19")
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

// Generate generates all code.
func Generate() error {
	return common.GenerateWith([]string{
		filepath.Dir(deps.GoBinPath("mockgen")),
		filepath.Dir(deps.GoBinPath("wire")),
	})
}

// Cleans the bin director
func Clean() error {
	return os.RemoveAll("dist")
}

func Deps() {
	deps.GetAllDeps()
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

// All runs all targets in the appropriate order.
// The targets are run in the following order:
// deps, generate, lint, test, build, dockerImage
func All() error {
	mg.SerialDeps(Deps, Generate, Lint, Test)
	return nil
}
