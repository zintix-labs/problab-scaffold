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

// Package devgames holds the ProbLab Dev panel plugins, one <gid>.js per game,
// injected by cmd/dev. A game without a plugin is shown with the common panel.
// See server/api/dev/README.md in Problab for the plugin interface.
package devgames

import "embed"

// FS provides the panel plugins (<gid>.js) at its root.
//
//go:embed *.js
var FS embed.FS
