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

package tags

import (
	"github.com/zintix-labs/problab/sdk/tag"
)

// Catalog binds each game's collection and analysis tags to its spec.GID. This is the
// one place in the repository that owns this binding — per-game tag files
// (such as demo_0_tags.go) define the predicates. Both cmd/opt and cmd/exp
// inject this catalog as their default gameTags.
var Catalog = tag.GameTagCatalog{
	0: Demo_0_Tags,
}
