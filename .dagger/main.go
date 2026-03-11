// A generated module for ZeroDowntime functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/zero-downtime/internal/dagger"
)

type ZeroDowntime struct{}

// Image builds a binary (client or server) and packages it into a scratch container.
func (m *ZeroDowntime) BuildImage(
	ctx context.Context,
	// Source directory of the project
	// +defaultPath="/"
	source *dagger.Directory,
	// Protobuf definition directory
	// +defaultPath="/helloworld"
	protobuf *dagger.Directory,
	// Name of the binary to build, either "client" or "server"
	// +default="client"
	target string,
) *dagger.Container {

	generated := dag.Protobuf(dagger.ProtobufOpts{Source: protobuf}).
		Generate().
		Directory("/out/")
	
	binary := dag.Golang(dagger.GolangOpts{Source: source.WithDirectory("helloworld/", generated)}).
		WithCgoDisabled().
		Build(dagger.GolangBuildOpts{Args: []string{"./" + target}}).
		File(target)
	return dag.Container().
		WithFile("/"+target, binary).
		WithEntrypoint([]string{"/"+target})
}
