// Package converter holds the scaffold's application-owned dto.ResultConverter
// implementations. cmd/opt injects one of them through
// optimizerv2.WithResultConverter to encode each RGS results.jsonl.zst line.
package converter

import (
	"encoding/json"

	"github.com/zintix-labs/problab/dto"
)

// IdentityConverter is the scaffold's starting-point converter. It mirrors
// dto.IdentityConverter and emits the whole dto.SpinResult, including the start
// and after PRNG snapshots and checkpoint: do not disclose its output to parties
// that should not receive the underlying replay material.
//
// Because it is a distinct function from dto.IdentityConverter, the optimizer/v2
// run report labels it converter=custom even though the payload is identical.
// Replace it (or add a sibling) with a converter that encodes your platform's
// result schema and drops the replay state before shipping RGS output.
func IdentityConverter(r dto.SpinResult) (json.RawMessage, error) {
	return json.Marshal(r)
}
