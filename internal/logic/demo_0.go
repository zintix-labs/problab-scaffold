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

package logic

import (
	"log"

	"github.com/zintix-labs/problab/sdk/buf"
	"github.com/zintix-labs/problab/sdk/slot"
	"github.com/zintix-labs/problab/spec"
)

// ============================================================
// ** Registration **
// ============================================================

func init() {
	logic := "demo_normal"
	if err := slot.GameRegister[*buf.NoExtend](
		spec.LogicKey(logic),
		buildGame0000,
		Logics,
	); err != nil {
		log.Fatalf("%s register failed: %v", logic, err)
	}
}

// ============================================================
// ** Game Interface **
// ============================================================

type game0000 struct {
	fixed *fixed0000
	ext   *ext0000
}

func buildGame0000(gh *slot.Game) (slot.GameLogic, error) {
	g := &game0000{
		fixed: new(fixed0000),
		ext:   nil,
	}
	if err := spec.DecodeFixed(gh.GameSetting, g.fixed); err != nil {
		return nil, err
	}
	g.fixed.symboltypes = gh.GameSetting.GameModeSettings[0].SymbolSetting.SymbolTypes
	g.ext = g.newext(gh.GameSetting.GameModeSettings[0].ScreenSetting.ScreenSize, gh.IsSim)
	return g, nil
}

// ============================================================
// ** Game-specific Fixed Configuration **
// ============================================================

// fixed
type fixed0000 struct {
	FreeRound   int    `yaml:"free_round"`
	DemoB       []int  `yaml:"demo_b"`
	DemoC       string `yaml:"demo_c"`
	symboltypes []spec.SymbolType
}

// ============================================================
// ** Game-specific Extension State (implements Reset and Snapshot) **
// ============================================================

type ext0000 struct {
	Triggered     bool  `json:"is_trigger"`
	ScatterCount  int   `json:"scatters,omitzero"`
	ScatterHitMap []int `json:"scatter_hits,omitzero"`
	isSim         bool
}

func (g *game0000) newext(screensize int, isSim bool) *ext0000 {
	return &ext0000{
		Triggered:     false,
		ScatterCount:  0,
		ScatterHitMap: make([]int, 0, screensize),
		isSim:         isSim,
	}
}

func (e *ext0000) Reset() {
	e.Triggered = false
	e.ScatterCount = 0
	e.ScatterHitMap = e.ScatterHitMap[:0]
}

func (e *ext0000) Snapshot() any {
	if e.isSim {
		return nil
	}
	hits := make([]int, len(e.ScatterHitMap))
	copy(hits, e.ScatterHitMap)
	ec := &ext0000{
		Triggered:     e.Triggered,
		ScatterCount:  e.ScatterCount,
		ScatterHitMap: hits,
	}
	return ec
}

// ============================================================
// ** Main Game Logic Entry **
// ============================================================

// GetResult is the main entry point and returns the final *SpinResult
func (g *game0000) GetResult(r *buf.SpinRequest, gh *slot.Game) *buf.SpinResult {
	sr := gh.StartNewSpin(r)

	base := g.getBaseResult(r.BetMult, gh)
	sr.AppendModeResult(base)

	if base.Trigger != 0 {
		free := g.getFreeResult(r.BetMult, gh)
		sr.AppendModeResult(free)
	}
	sr.End()
	return sr
}

// ============================================================
// ** Per-Mode Internal Logic Implementation **
// ============================================================

func (g *game0000) getBaseResult(betMult int, gh *slot.Game) *buf.GameModeResult {
	mode := gh.GameModeHandlerList[0]
	sg := mode.ScreenGenerator
	sc := mode.ScreenCalculator
	gmr := mode.GameModeResult
	ext := g.ext
	ext.Reset()

	// 1. Generate screen
	screen := sg.GenScreen()
	gmr.AddAct(buf.FinishAct, "screen", screen, nil)

	// 2. Calculate win
	sc.CalcScreen(betMult, screen, gmr)
	if gmr.GetTmpWin() > 0 {
		gmr.AddAct(buf.FinishAct, "win", nil, nil)
	}

	// 3. Check trigger condition
	gmr.Trigger = g.trigger(screen)
	if gmr.Trigger > 0 {
		gmr.AddAct(buf.FinishAct, "trigger", nil, ext)
	}

	// 4. Commit round result
	gmr.FinishRound()

	return mode.YieldResult()
}

func (g *game0000) getFreeResult(betMult int, gh *slot.Game) *buf.GameModeResult {
	mode := gh.GameModeHandlerList[1]
	sg := mode.ScreenGenerator
	sc := mode.ScreenCalculator
	gmr := mode.GameModeResult
	round := g.fixed.FreeRound
	ext := g.ext
	ext.Reset()

	for i := 0; i < round; i++ {
		// 1. Generate screen
		screen := sg.GenScreen()
		gmr.AddAct(buf.FinishAct, "screen", screen, nil)

		// 2. Calculate win
		sc.CalcScreen(betMult, screen, gmr)
		if gmr.GetTmpWin() > 0 {
			gmr.AddAct(buf.FinishAct, "win", nil, nil)
		}

		// 3. Finalize round actions
		gmr.FinishRound()
	}

	return mode.YieldResult()
}

// ============================================================
// ** Internal Helper Functions **
// ============================================================

// Returns 0 if not triggered; >0 indicates trigger
func (g *game0000) trigger(screen []int16) int {
	g.ext.Reset()
	ext := g.ext
	symtype := g.fixed.symboltypes
	for i := 0; i < len(screen); i++ {
		if symtype[screen[i]] == spec.SymbolTypeScatter {
			ext.ScatterCount++
			ext.ScatterHitMap = append(ext.ScatterHitMap, i)
		}
	}
	if ext.ScatterCount > 2 {
		ext.Triggered = true
		return 1
	}
	return 0
}
