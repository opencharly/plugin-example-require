// Package examplerequire is the reference charly plugin for the DECLARED inter-plugin
// dependency mechanism (`plugin.requires:`).
//
// Its candy manifest (candy/plugin-example-require/charly.yml) declares TWO peers —
// `verb:externalprobe` (served by candy/plugin-example-external, compiled into charly)
// and `verb:exampledispatchpeer` (served out-of-process by candy/plugin-example-dispatch)
// — and its Invoke PEER-INVOKES both over the host broker
// (sdk.Executor.InvokeProvider), echoing their results so a bed can assert the
// declaration resolved BOTH placements (already-registered + on-demand connect).
//
// The authored manifest is THE source of the declaration: the host reads it from the
// resolved view (spec.CandyReader.GetPluginRequires) and resolves every requirement
// before this plugin's providers register. That is why there is deliberately NO
// sdk.NewMetaWithRequires call here — the author declares a peer ONCE, in the manifest,
// with no Go change and no charly-core change per new plugin.
package examplerequire

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"

	"github.com/opencharly/sdk"
	pb "github.com/opencharly/spec/proto"
)

//go:embed schema/*.cue
var schemaFS embed.FS

const calver = "2026.269.0001"

// NewProvider returns the examplerequire provider.
func NewProvider() pb.ProviderServer { return &provider{} }

// NewMeta advertises verb:examplerequire + the plugin's self-contained CUE input schema
// (via sdk.NewMeta → BuildCapabilities, which compiles the schema standalone, failing
// loudly if broken). The declared PEER requirements live in the candy manifest
// (`plugin.requires:`), NOT here — the manifest is THE authored source the host reads.
func NewMeta() pb.PluginMetaServer {
	return sdk.NewMeta(calver,
		[]sdk.ProvidedCapability{{Class: "verb", Word: "examplerequire", InputDef: "#ExamplerequireInput"}},
		schemaFS)
}

type provider struct{ pb.UnimplementedProviderServer }

// requireInput is the verb's plugin_input: a marker plus the two peer verb words to
// invoke. The peer words arrive as DATA so this plugin's Go never restates its manifest
// declaration (no manifest/Go duplication) — the manifest declares the dependency, the
// caller names the peers, and the load gate has already resolved them.
type requireInput struct {
	Marker       string `json:"marker"`
	BuiltinPeer  string `json:"builtin_peer"`  // the compiled-in peer (resolve branch)
	ExternalPeer string `json:"external_peer"` // the out-of-process peer (connect branch)
}

// Invoke peer-invokes each named peer over the host broker and returns their results,
// keyed by branch so a bed can assert both. A peer the requires gate could not resolve
// would already have failed the load loudly; a peer that vanished after the load surfaces
// here as an InvokeProvider error.
func (provider) Invoke(ctx context.Context, req *pb.InvokeRequest) (*pb.InvokeReply, error) {
	if req.GetOp() != sdk.OpRun {
		return nil, fmt.Errorf("examplerequire: unsupported op %q (only %q)", req.GetOp(), sdk.OpRun)
	}
	exec, err := sdk.ExecutorFromInvoke(req.GetExecutorBrokerId())
	if err != nil {
		return nil, fmt.Errorf("examplerequire: no host executor: %w", err)
	}
	var in requireInput
	if len(req.GetParamsJson()) > 0 {
		if err := json.Unmarshal(req.GetParamsJson(), &in); err != nil {
			return nil, fmt.Errorf("examplerequire: decode input: %w", err)
		}
	}
	out := map[string]any{"status": "pass", "message": in.Marker}
	// A declared verb peer is invoked with its own plugin_input; the reference peers read
	// a `marker` (externalprobe) or an optional peer command (exampledispatchpeer, whose
	// absent value is the plain peer-reached echo).
	params := []byte(`{"plugin_input":{"marker":"requires-peer"}}`)
	for label, word := range map[string]string{"builtin_peer": in.BuiltinPeer, "external_peer": in.ExternalPeer} {
		if word == "" {
			continue
		}
		pres, ierr := exec.InvokeProvider(ctx, "verb", word, sdk.OpRun, params, nil, sdk.InvokeProviderOpts{})
		if ierr != nil {
			return nil, fmt.Errorf("examplerequire: invoke declared peer %q: %w", word, ierr)
		}
		out[label] = json.RawMessage(pres)
	}
	res, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}
	return &pb.InvokeReply{ResultJson: res}, nil
}
