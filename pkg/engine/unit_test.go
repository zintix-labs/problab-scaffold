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

package engine

import (
	"io/fs"
	"testing"

	"github.com/zintix-labs/problab/catalog"
	"github.com/zintix-labs/problab/spec"
)

func TestConfigsEmbedded(t *testing.T) {
	if len(cfgs) == 0 {
		t.Fatal("cfgs is empty")
	}
	for _, name := range []string{"demo_0.yaml", "demo_1.yaml"} {
		if !configExists(name) {
			t.Fatalf("config not found in embedded FS: %s", name)
		}
	}
}

func TestLogicRegistry(t *testing.T) {
	if len(logics) == 0 {
		t.Fatal("logics is empty")
	}
	reg := logics[0]
	if reg == nil {
		t.Fatal("logic registry is nil")
	}
	if !reg.IsExist(spec.LogicKey("demo_normal")) {
		t.Error("missing logic key: demo_normal")
	}
	if !reg.IsExist(spec.LogicKey("demo_cascade")) {
		t.Error("missing logic key: demo_cascade")
	}
}

func TestNewAndSummary(t *testing.T) {
	lab, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if lab == nil {
		t.Fatal("New() returned nil")
	}

	ent0, ok := lab.EntryById(spec.GID(0))
	if !ok {
		t.Fatal("missing entry for game id 0")
	}
	if ent0.ConfigName != "demo_0.yaml" {
		t.Fatalf("unexpected config for game id 0: %s", ent0.ConfigName)
	}

	ent1, ok := lab.EntryById(spec.GID(1))
	if !ok {
		t.Fatal("missing entry for game id 1")
	}
	if ent1.ConfigName != "demo_1.yaml" {
		t.Fatalf("unexpected config for game id 1: %s", ent1.ConfigName)
	}

	sum, err := lab.Summary()
	if err != nil {
		t.Fatalf("Summary() error: %v", err)
	}
	if len(sum) < 2 {
		t.Fatalf("expected at least 2 summaries, got %d", len(sum))
	}
	byID := make(map[spec.GID]catalog.Summary, len(sum))
	for _, s := range sum {
		byID[s.GID] = s
	}
	if s, ok := byID[spec.GID(0)]; !ok || s.Logic != spec.LogicKey("demo_normal") {
		t.Fatal("summary missing or invalid for game id 0")
	}
	if s, ok := byID[spec.GID(1)]; !ok || s.Logic != spec.LogicKey("demo_cascade") {
		t.Fatal("summary missing or invalid for game id 1")
	}
}

func TestNewMachine(t *testing.T) {
	lab, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	for _, id := range []spec.GID{0, 1} {
		m, err := lab.NewMachine(id, true)
		if err != nil {
			t.Fatalf("NewMachine(%d) error: %v", id, err)
		}
		if m == nil {
			t.Fatalf("NewMachine(%d) returned nil", id)
		}
	}
}

func TestNewValidation(t *testing.T) {
	origFactory := pRNGFactory
	origCfgs := cfgs
	origLogics := logics
	t.Cleanup(func() {
		pRNGFactory = origFactory
		cfgs = origCfgs
		logics = origLogics
	})

	reset := func() {
		pRNGFactory = origFactory
		cfgs = origCfgs
		logics = origLogics
	}

	t.Run("nil PRNG factory", func(t *testing.T) {
		reset()
		pRNGFactory = nil
		if _, err := New(); err == nil {
			t.Fatal("expected error for nil PRNG factory")
		}
	})

	t.Run("missing configs", func(t *testing.T) {
		reset()
		cfgs = nil
		if _, err := New(); err == nil {
			t.Fatal("expected error for missing configs")
		}
	})

	t.Run("missing logics", func(t *testing.T) {
		reset()
		logics = nil
		if _, err := New(); err == nil {
			t.Fatal("expected error for missing logics")
		}
	})
}

func configExists(name string) bool {
	for _, src := range cfgs {
		if src == nil {
			continue
		}
		if _, err := fs.ReadFile(src, name); err == nil {
			return true
		}
	}
	return false
}
