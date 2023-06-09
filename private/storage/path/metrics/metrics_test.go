// Copyright 2019 Anapaya Systems
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

package pathdb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scionproto/scion/private/storage"
	pathdbtest "github.com/scionproto/scion/private/storage/path/dbtest"
	pathstoragemetrics "github.com/scionproto/scion/private/storage/path/metrics"
	"github.com/scionproto/scion/private/storage/path/sqlite"
)

type TestPathDB struct {
	storage.PathDB
}

func (b *TestPathDB) Prepare(t *testing.T, _ context.Context) {
	b.PathDB = setupDB(t)
}

// TestMetricWrapperFunctionality tests that the metrics wrapper succeeds the
// normal path db test suite.
func TestMetricWrapperFunctionality(t *testing.T) {
	tdb := &TestPathDB{}
	pathdbtest.TestPathDB(t, tdb)
}

func setupDB(t *testing.T) storage.PathDB {
	db, err := sqlite.New("file::memory:")
	require.NoError(t, err)
	return pathstoragemetrics.WrapDB(db, pathstoragemetrics.Config{
		Driver: "test_sqlite",
	})
}
