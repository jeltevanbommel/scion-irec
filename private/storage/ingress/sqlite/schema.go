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

package sqlite

const (
	// SchemaVersion is the version of the SQLite schema understood by this backend.
	// Whenever changes to the schema are made, this version number should be increased
	// to prevent data corruption between incompatible database schemas.
	SchemaVersion = 1
	// Schema is the SQLite database layout.
	Schema = `CREATE TABLE Beacons(
		RowID INTEGER PRIMARY KEY,
		SegID DATA NOT NULL,
		FullID DATA UNIQUE NOT NULL,
		StartIsd INTEGER NOT NULL,
		StartAs INTEGER NOT NULL,
		StartIntfGroup INTEGER NOT NULL,
		AlgorithmHash DATA NOT NULL,
		AlgorithmId DATA NOT NULL,
		InIntfID INTEGER NOT NULL,
		HopsLength INTEGER NOT NULL,
		InfoTime INTEGER NOT NULL,
		ExpirationTime INTEGER NOT NULL,
		LastUpdated INTEGER NOT NULL,
		Usage INTEGER NOT NULL,
		Marker INTEGER NOT NULL,
		Beacon BLOB NOT NULL
	);


CREATE TABLE Algorithm(
		RowID INTEGER PRIMARY KEY,
		AlgorithmHash DATA UNIQUE NOT NULL,
		Algorithm BLOB NOT NULL,
		Status INTEGER NOT NULL,
		FetchTime INTEGER NOT NULL,
		ExpirationTime INTEGER NOT NULL,
		LastUsed INTEGER NOT NULL,
		Usage INTEGER NOT NULL,
		Beacon BLOB NOT NULL
	);
	`
	// marker 0 = new
	// marker 1 = already processed
	// Algorithm statuses:
	// 0: Unfetched
	// 1: Fetching -> check FetchTime whether expired ?
	// 2: Installed
	// 3: Manually installed (e.g. from disk)
	BeaconsTable = "Beacons"
)
