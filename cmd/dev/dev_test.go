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

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/zintix-labs/problab"
)

// The default injection points: pLab builds the project's games and devGames
// provides a <gid>.js for each of them.
func TestDefaultInjectionPoints(t *testing.T) {
	// The engine resolves internal/optimal from the project root.
	t.Chdir("../..")
	lab, err := pLab()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = lab.Close() })
	ids := lab.IDs()
	if len(ids) == 0 {
		t.Fatal("no games registered")
	}
	for _, gid := range ids {
		raw, err := fs.ReadFile(devGames, fmt.Sprintf("%d.js", gid))
		if err != nil {
			t.Fatalf("game %d has no plugin: %v", gid, err)
		}
		if !strings.Contains(string(raw), fmt.Sprintf("DEV.register(%d,", gid)) {
			t.Fatalf("%d.js does not register game %d", gid, gid)
		}
	}
}

// Injection errors are reported; unknown arguments are rejected before the Lab
// is built.
func TestInjectionAndArguments(t *testing.T) {
	old := pLab
	t.Cleanup(func() { pLab = old })
	sentinel := errors.New("private lab failure")
	pLab = func() (*problab.Problab, error) { return nil, sentinel }
	if code, err := run(nil, nil); code != 1 || !errors.Is(err, sentinel) {
		t.Fatal(code, err)
	}
	pLab = func() (*problab.Problab, error) { t.Fatal("Lab built before arguments were checked"); return nil, nil }
	for _, args := range [][]string{{"-config", "x"}, {"extra"}} {
		if code, err := run(args, nil); code != 1 || err == nil {
			t.Fatal(args, code, err)
		}
	}
}
