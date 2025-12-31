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

package configs

import (
	"embed"
)

// FS provides embedded default config YAMLs for this scaffold.
//
// Important:
//   - Do NOT delete or move this file unless you also update bootstrap wiring.
//   - Problab-scaffold expects a *flat* config FS layout by default:
//   - YAML files are addressed by filename only (no nested paths)
//   - configs are embedded from the same folder as this file
//
// If you want to organize configs into subfolders, you must also adjust:
//   - the embed pattern(s) below
//   - the bootstrap/config loading logic that consumes this FS
//
// Embed all YAML files in this directory (flat layout).
//
//go:embed *.yaml
var FS embed.FS
