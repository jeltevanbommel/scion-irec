// Copyright 2018 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scionproto/scion/pkg/private/xtest/graph"
)

func TestGeneratedUpToDate(t *testing.T) {
	g, err := LoadGraph("../../../../" + DefaultTopoFile)
	require.NoError(t, err)
	graphMapping := make(map[string]int)
	for i, id := range g.IfaceIds {
		graphMapping[i.Name()] = id
	}
	assert.Equal(t, graph.StaticIfaceIdMapping, graphMapping,
		"Generated graph is out of date, run graphupdater")
}
