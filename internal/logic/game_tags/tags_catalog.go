// Copyright 2025 Zintix Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package game_tags

import (
	"github.com/zintix-labs/problab/sdk/tag"
	"github.com/zintix-labs/problab/spec"
)

// GameTagCatalog maps each game's spec.GID to its named collection tags. It has
// the same underlying type that optimizerv2.WithCollectionTags accepts, so a
// catalog value is passed to the Tuner without conversion. Tag names are the
// ones opt_cfg.yaml references from classes[].collect.tags; bg and fg are
// built in and must not be redefined.
type GameTagCatalog map[spec.GID]map[string]tag.IsTag

// GameTags binds each game's collection tag set to its spec.GID. This is the
// one place in the repository that owns this binding — per-game tag files
// (such as demo_0_tags.go) define the predicates, and cmd/opt injects this
// catalog as its default gameTags.
var GameTags = GameTagCatalog{
	0: Demo_0_Tags,
}
