package examplerequire

import (
	"context"
	"testing"

	pb "github.com/opencharly/spec/proto"
)

// The plugin's declared peer dependencies live in the candy manifest
// (`plugin.requires:`), NOT in its Go — there is deliberately no
// sdk.NewMetaWithRequires call. This locks that contract: Describe advertises the
// provided verb with an EMPTY wire requires, so the ONE authored source of the
// declaration is the manifest, read by the host from the resolved view.
func TestNewMeta_AdvertisesProviderAndNoWireRequires(t *testing.T) {
	caps, err := NewMeta().Describe(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("Describe: %v", err)
	}
	found := false
	for _, c := range caps.GetProvided() {
		if c.GetClass() == "verb" && c.GetWord() == "examplerequire" {
			found = true
		}
	}
	if !found {
		t.Fatal("NewMeta must advertise verb:examplerequire")
	}
	if len(caps.GetRequires()) != 0 {
		t.Fatalf("wire requires must be empty (the manifest is THE declaration), got %d", len(caps.GetRequires()))
	}
}
