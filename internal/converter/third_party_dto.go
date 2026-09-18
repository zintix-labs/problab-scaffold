// Copyright 2025 Zintix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package converter holds the scaffold's application-owned dto.ResultConverter
// implementations. cmd/opt injects one of them through
// optimizerv2.WithResultConverter to encode each RGS results.jsonl.zst line.
package converter

import (
	"encoding/json"

	"github.com/zintix-labs/problab/dto"
	"github.com/zintix-labs/problab/spec"
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

// SpinResultWithoutSnapshot is the payload handed to an RGS or platform: the
// same round data as dto.SpinResult, with the engine's Core snapshots removed.
//
// The field list is deliberately explicit rather than embedding dto.SpinResult.
// A field added to SpinResult upstream must not start flowing into an exported
// bank on its own; with this shape the default is to leave it out.
//
// This removes only what problab put in State. A game's own ExtendResult
// payload (gamemodes[].acts[].ext) is authored by the game and is passed
// through verbatim — auditing that stays with the game author.
type SpinResultWithoutSnapshot struct {
	GameName   string                  `json:"game"`
	GameID     spec.GID                `json:"gameid"`
	TotalWin   int                     `json:"win"`
	Bet        int                     `json:"bet"`
	BetMode    int                     `json:"betmode"`
	BetMult    int                     `json:"betmult"`
	GameModes  []dto.GameModeResultDTO `json:"gamemodes,omitempty"`
	IsGameEnd  bool                    `json:"isend"`
	Checkpoint json.RawMessage         `json:"cp,omitempty"` // 完整局結束的遊戲可整欄拿掉
}

// WithoutSnapshotConverter is a dto.ResultConverter. The DTO is borrowed for the duration of the
// call; every field below is copied or re-referenced, nothing is mutated.
func WithoutSnapshotConverter(r dto.SpinResult) (json.RawMessage, error) {
	return json.Marshal(SpinResultWithoutSnapshot{
		GameName:   r.GameName,
		GameID:     r.GameID,
		TotalWin:   r.TotalWin,
		Bet:        r.Bet,
		BetMode:    r.BetMode,
		BetMult:    r.BetMult,
		GameModes:  r.GameModes,
		IsGameEnd:  r.IsGameEnd,
		Checkpoint: r.State.Checkpoint,
	})
}
