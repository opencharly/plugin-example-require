// Self-contained input schema for the verb:examplerequire capability — references no
// base def, so it compiles standalone (gengotypes + the SDK's serve-side compile).
#ExamplerequireInput: {
	marker?:        string
	builtin_peer?:  string
	external_peer?: string
}
