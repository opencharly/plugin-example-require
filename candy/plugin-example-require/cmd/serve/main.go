// Command serve is the OUT-OF-PROCESS entrypoint for the examplerequire reference
// plugin: a thin shim serving the importable provider over go-plugin gRPC via
// sdk.Serve. The SAME NewProvider()/NewMeta() compile INTO charly in-process when the
// candy is listed in compiled_plugins; this binary is host-built + connected only when
// it is not — placement is invisible above the registry.
package main

import (
	examplerequire "github.com/opencharly/plugin-example-require/candy/plugin-example-require"
	"github.com/opencharly/sdk"
)

func main() { sdk.Serve(examplerequire.NewProvider(), examplerequire.NewMeta()) }
