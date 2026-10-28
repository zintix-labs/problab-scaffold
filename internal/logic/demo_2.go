// Copyright 2026 Zintix Labs
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

package logic

import (
	"log"

	"github.com/zintix-labs/problab/dto"
	"github.com/zintix-labs/problab/sdk/buf"
	"github.com/zintix-labs/problab/sdk/slot"
	"github.com/zintix-labs/problab/spec"
)

// demo_choice: a segmented game example. When BG triggers, segment 0 stops at a
// choice point and the result's ChooseNext lists two entries; after the player
// chooses, segment 1 runs the chosen FG and ends the game.

// ============================================================
// ** Registration **
// ============================================================

func init() {
	logic := "demo_choice"
	if err := slot.GameRegister(spec.LogicKey(logic), buildGame0002, Logics); err != nil {
		log.Fatalf("%s register failed: %v", logic, err)
	}
	if err := dto.RegisterExtendRender[ext0002](spec.LogicKey(logic)); err != nil {
		log.Fatalf("%s register failed: %v", logic, err)
	}
	// Continuation segments receive their entry state through a checkpoint:
	// register its codec so the production path can round-trip it through JSON.
	if err := dto.RegisterCheckpoint[check0002](spec.LogicKey(logic)); err != nil {
		log.Fatalf("%s register failed: %v", logic, err)
	}
}

// Entry values (the contents of ChooseNext). They deliberately differ from the
// position indexes: a continuation segment reads its entry as ChooseNext[Choice].
const (
	Demo0002EntryFGA = 1 // more free spins
	Demo0002EntryFGB = 2 // fewer free spins, multiplied wins
)

// ============================================================
// ** Game Interface **
// ============================================================

type game0002 struct {
	fixed *fixed0002
	ext   *ext0002
}

func buildGame0002(gh *slot.Game) (slot.GameLogic, error) {
	g := &game0002{fixed: new(fixed0002)}
	if err := spec.DecodeFixed(gh.GameSetting, g.fixed); err != nil {
		return nil, err
	}
	g.fixed.symboltypes = gh.GameSetting.GameModeSettings[0].SymbolSetting.SymbolTypes
	screensize := gh.GameSetting.GameModeSettings[0].ScreenSetting.ScreenSize
	g.ext = &ext0002{ScatterHitMap: make([]int, 0, screensize), isSim: gh.OmitExtend()}
	return g, nil
}

// ============================================================
// ** Game-specific Fixed Configuration **
// ============================================================

type fixed0002 struct {
	FreeRoundA  int `yaml:"free_round_a"`
	FreeRoundB  int `yaml:"free_round_b"`
	MultiplierB int `yaml:"multiplier_b"`
	symboltypes []spec.SymbolType
}

// ============================================================
// ** Extend and Checkpoint **
// ============================================================

type ext0002 struct {
	Triggered     bool  `json:"is_trigger"`
	ScatterCount  int   `json:"scatters,omitzero"`
	ScatterHitMap []int `json:"scatter_hits,omitzero"`
	isSim         bool
}

func (e *ext0002) Reset() {
	e.Triggered = false
	e.ScatterCount = 0
	e.ScatterHitMap = e.ScatterHitMap[:0]
}

func (e *ext0002) Snapshot() any {
	if e.isSim {
		return nil
	}
	return &ext0002{
		Triggered:     e.Triggered,
		ScatterCount:  e.ScatterCount,
		ScatterHitMap: append([]int(nil), e.ScatterHitMap...),
	}
}

// check0002 is the entry state segment 0 hands to the continuation segment.
// It is newly allocated each time and never modified after it is handed over.
type check0002 struct {
	Scatters int `json:"scatters"`
}

// ============================================================
// ** Main Game Logic Entry **
// ============================================================

// GetResult runs BG in segment 0 and stops at a choice point when it triggers.
// A continuation segment first validates its entry and checkpoint (returning
// nil without consuming RNG when they are invalid), then runs the chosen FG.
func (g *game0002) GetResult(r *buf.SpinRequest, gh *slot.Game) *buf.SpinResult {
	sr := gh.StartNewSpin(r)
	if r.Cycle == 0 {
		base, scatters := g.getBaseResult(r, gh)
		sr.AppendModeResult(base)
		if base.Trigger > 0 {
			sr.ChooseNext = append(sr.ChooseNext, Demo0002EntryFGA, Demo0002EntryFGB)
			sr.State.Checkpoint = &check0002{Scatters: scatters}
		} else {
			sr.End()
		}
		return sr
	}

	var fg *buf.GameModeResult
	switch r.ChooseNext[r.Choice] {
	case Demo0002EntryFGA:
		fg = g.getFreeAResult(r, gh)
	case Demo0002EntryFGB:
		fg = g.getFreeBResult(r, gh)
	default:
		return nil
	}
	if fg == nil {
		return nil
	}

	sr.AppendModeResult(fg)
	sr.End()
	return sr
}

// ============================================================
// ** Per-Mode Internal Logic Implementation **
// ============================================================

func (g *game0002) getBaseResult(r *buf.SpinRequest, gh *slot.Game) (*buf.GameModeResult, int) {
	mode := gh.GameModeHandlerList[0]
	gmr := mode.GameModeResult
	ext := g.ext
	ext.Reset()

	screen := mode.ScreenGenerator.GenScreen()
	gmr.AddAct(buf.FinishAct, "screen", screen, nil)

	mode.ScreenCalculator.CalcScreen(r.BetMult, screen, gmr)
	if gmr.GetTmpWin() > 0 {
		gmr.AddAct(buf.FinishAct, "win", nil, nil)
	}

	scatters := g.scatters(screen)
	if scatters >= 3 {
		ext.Triggered = true
		gmr.Trigger = 1
		gmr.AddAct(buf.FinishAct, "trigger", nil, ext)
	}
	gmr.FinishRound()
	return mode.YieldResult(), scatters
}

func (g *game0002) getFreeAResult(r *buf.SpinRequest, gh *slot.Game) *buf.GameModeResult {
	cp, ok := g.parseCheckpoint(r)
	if !ok {
		return nil
	}
	rounds := g.fixed.FreeRoundA + cp.Scatters - 3 // one extra spin per extra Scatter
	return g.getFreeResult(1, rounds, r.BetMult, gh)
}

func (g *game0002) getFreeBResult(r *buf.SpinRequest, gh *slot.Game) *buf.GameModeResult {
	cp, ok := g.parseCheckpoint(r)
	if !ok {
		return nil
	}
	rounds := g.fixed.FreeRoundB + cp.Scatters - 3
	betMult := r.BetMult * g.fixed.MultiplierB
	return g.getFreeResult(2, rounds, betMult, gh)
}

// parseCheckpoint validates the continuation segment's entry state before any
// screen is generated.
func (g *game0002) parseCheckpoint(r *buf.SpinRequest) (*check0002, bool) {
	cp, ok := r.Checkpoint.(*check0002)
	if !ok || cp.Scatters < 3 {
		return nil, false
	}
	return cp, true
}

func (g *game0002) getFreeResult(modeIndex, rounds, betMult int, gh *slot.Game) *buf.GameModeResult {
	mode := gh.GameModeHandlerList[modeIndex]
	gmr := mode.GameModeResult
	for range rounds {
		screen := mode.ScreenGenerator.GenScreen()
		gmr.AddAct(buf.FinishAct, "screen", screen, nil)
		mode.ScreenCalculator.CalcScreen(betMult, screen, gmr)
		if gmr.GetTmpWin() > 0 {
			gmr.AddAct(buf.FinishAct, "win", nil, nil)
		}
		gmr.FinishRound()
	}
	return mode.YieldResult()
}

// scatters counts the Scatter symbols on the screen and records them in ext.
func (g *game0002) scatters(screen []int16) int {
	ext := g.ext
	for i, sym := range screen {
		if g.fixed.symboltypes[sym] == spec.SymbolTypeScatter {
			ext.ScatterCount++
			ext.ScatterHitMap = append(ext.ScatterHitMap, i)
		}
	}
	return ext.ScatterCount
}
