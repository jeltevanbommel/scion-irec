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

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/serrors"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/scionproto/scion/control/beacon"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/private/storage/db"
)

//var _ beacon.DB = (*Backend)(nil)

type Backend struct {
	db *sql.DB
	*executor
}

// New returns a new SQLite backend opening a database at the given path. If
// no database exists a new database is be created. If the schema version of the
// stored database is different from the one in schema.go, an error is returned.
func New(path string, ia addr.IA) (*Backend, error) {
	db, err := db.NewSqlite(path, Schema, SchemaVersion)
	if err != nil {
		return nil, err
	}
	return &Backend{
		executor: &executor{
			db: db,
			ia: ia,
		},
		db: db,
	}, nil
}

// SetMaxOpenConns sets the maximum number of open connections.
func (b *Backend) SetMaxOpenConns(maxOpenConns int) {
	b.db.SetMaxOpenConns(maxOpenConns)
}

// SetMaxIdleConns sets the maximum number of idle connections.
func (b *Backend) SetMaxIdleConns(maxIdleConns int) {
	b.db.SetMaxIdleConns(maxIdleConns)
}

// Close closes the database.
func (b *Backend) Close() error {
	return b.db.Close()
}

type executor struct {
	sync.RWMutex
	db db.Sqler
	ia addr.IA
}

type beaconMeta struct {
	RowID       int64
	InfoTime    time.Time
	LastUpdated time.Time
}

func (e *executor) BeaconSources(ctx context.Context, ignoreIntfGroup bool) ([]*cppb.RACBeaconSource, error) {
	e.RLock()
	defer e.RUnlock()
	var query string
	if ignoreIntfGroup {
		query = `SELECT DISTINCT StartIsd, StartAs, 0, AlgorithmHash, AlgorithmId FROM Beacons WHERE Marker = 0`
	} else {
		query = `SELECT DISTINCT StartIsd, StartAs, StartIntfGroup, AlgorithmHash, AlgorithmId FROM Beacons WHERE Marker = 0`
	}
	rows, err := e.db.QueryContext(ctx, query)
	log.Info("query text", "text", query)
	if err != nil {
		return nil, db.NewReadError("Error selecting source IAs", err)
	}
	defer rows.Close()
	var beaconSources []*cppb.RACBeaconSource
	for rows.Next() {
		var isd addr.ISD
		var as addr.AS
		var algHash sql.RawBytes
		var algId uint32
		var intfGroup uint16

		if err := rows.Scan(&isd, &as, &intfGroup, &algHash, &algId); err != nil {
			return nil, err
		}

		ia, err := addr.IAFrom(isd, as)
		if err != nil {
			return nil, err
		}

		beaconSources = append(beaconSources, &cppb.RACBeaconSource{
			AlgorithmHash:   algHash,
			AlgorithmID:     algId,
			OriginAS:        uint64(ia),
			OriginIntfGroup: uint32(intfGroup),
		})
	}
	log.Info("query rows", "res", beaconSources)
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return beaconSources, nil
}

// InsertBeacon inserts the beacon if it is new or updates the changed
// information.
func (e *executor) InsertBeacon(
	ctx context.Context,
	b beacon.Beacon,
	usage beacon.Usage,
) (beacon.InsertStats, error) {

	ret := beacon.InsertStats{}
	// Compute ids outside of the lock.
	fullID := b.Segment.IRECID()

	e.Lock()
	defer e.Unlock()

	meta, err := e.getBeaconMeta(ctx, fullID)
	if err != nil {
		return ret, err
	}
	if meta != nil {
		// Update the beacon data if it is newer.
		if b.Segment.Info.Timestamp.After(meta.InfoTime) {
			if err := e.updateExistingBeacon(ctx, b, usage, meta.RowID, time.Now()); err != nil {
				return ret, err
			}
			ret.Updated = 1
			return ret, nil
		}
		return ret, nil
	}
	// Insert new beacon.
	err = db.DoInTx(ctx, e.db, func(ctx context.Context, tx *sql.Tx) error {
		return insertNewBeacon(ctx, tx, b, usage, time.Now())
	})
	if err != nil {
		return ret, err
	}

	ret.Inserted = 1
	return ret, nil

}

func (e *executor) GetAndMarkBeacons(ctx context.Context, maximum uint32, algHash []byte, algId uint32, originAS addr.IA, originIntfGroup uint32, ignoreIntfGroup bool, marker uint32) ([]*cppb.IRECBeacon, error) {
	e.Lock() // TODO(jvb): DB transaction instead of lock
	defer e.Unlock()
	selStmt, selArgs, updStmt, updArgs := e.buildSelUpdQuery(ctx, maximum, algHash, algId, originAS, originIntfGroup, ignoreIntfGroup, marker)
	rows, err := e.db.QueryContext(ctx, selStmt, selArgs...)
	if err != nil {
		return nil, serrors.WrapStr("looking up beacons", err, "query", selStmt)
	}
	defer rows.Close()
	var res []*cppb.IRECBeacon
	var count uint32 = 0
	for rows.Next() {
		var RowID int
		var lastUpdated int64
		var rawBeacon sql.RawBytes
		var InIntfID uint16
		var usage beacon.Usage
		err = rows.Scan(&RowID, &lastUpdated, &rawBeacon, &InIntfID, &usage)
		if err != nil {
			return nil, serrors.WrapStr("reading row", err)
		}
		seg, err := beacon.UnpackBeaconIREC(rawBeacon)

		if err != nil {
			return nil, serrors.WrapStr("parsing beacon", err)
		}
		res = append(res, &cppb.IRECBeacon{
			PathSeg: seg,
			InIfId:  uint32(InIntfID),
			Id:      count,
		})
		count += 1
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	_, err = e.db.ExecContext(ctx, updStmt, updArgs...)
	if err != nil {
		return nil, db.NewWriteError("updating markers", err)
	}
	return res, nil

}
func (e *executor) GetBeacons(ctx context.Context, maximum uint32, algHash []byte, algId uint32, originAS addr.IA, originIntfGroup uint32, ignoreIntfGroup bool) ([]*cppb.IRECBeacon, error) {

	e.RLock()
	defer e.RUnlock()
	stmt, args := e.buildQuery(ctx, maximum, algHash, algId, originAS, originIntfGroup, ignoreIntfGroup)
	rows, err := e.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, serrors.WrapStr("looking up beacons", err, "query", stmt)
	}
	defer rows.Close()
	var res []*cppb.IRECBeacon
	var count uint32 = 0
	for rows.Next() {
		var RowID int
		var lastUpdated int64
		var rawBeacon sql.RawBytes
		var InIntfID uint16
		var usage beacon.Usage
		err = rows.Scan(&RowID, &lastUpdated, &rawBeacon, &InIntfID, &usage)
		if err != nil {
			return nil, serrors.WrapStr("reading row", err)
		}
		seg, err := beacon.UnpackBeaconIREC(rawBeacon)

		if err != nil {
			return nil, serrors.WrapStr("parsing beacon", err)
		}
		res = append(res, &cppb.IRECBeacon{
			PathSeg: seg,
			InIfId:  uint32(InIntfID),
			Id:      count,
		})
		count += 1
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
func (e *executor) buildQuery(ctx context.Context, maximum uint32, algHash []byte, algId uint32, originAS addr.IA, originIntfGroup uint32, ignoreIntfGroup bool) (string, []interface{}) {
	var args []interface{}
	query := "SELECT DISTINCT RowID, LastUpdated, Beacon, InIntfID, Usage FROM Beacons"
	where := []string{}

	where = append(where, "(StartIsd=? AND StartAs=?)")
	args = append(args, originAS.ISD(), originAS.AS())
	//StartIntfGroup, AlgorithmHash
	where = append(where, "(AlgorithmHash=? AND AlgorithmID=?)")
	args = append(args, algHash, algId)
	if !ignoreIntfGroup {
		where = append(where, "(StartIntfGroup=?)")
		args = append(args, originIntfGroup)
	}
	if len(where) > 0 {
		query += "\n" + fmt.Sprintf("WHERE %s", strings.Join(where, " AND\n"))
	}

	if maximum > 0 {
		query += "\n LIMIT ?"
		args = append(args, maximum)
	}
	query += "\n" + "ORDER BY LastUpdated DESC"
	return query, args
}
func (e *executor) buildSelUpdQuery(ctx context.Context, maximum uint32, algHash []byte, algId uint32, originAS addr.IA, originIntfGroup uint32, ignoreIntfGroup bool, marker uint32) (string, []interface{}, string, []interface{}) {
	var selArgs []interface{}
	var updArgs []interface{}
	selQuery := "SELECT DISTINCT RowID, LastUpdated, Beacon, InIntfID, Usage FROM Beacons"
	updQuery := "UPDATE Beacons SET Marker = ?"
	updArgs = append(updArgs, marker)
	where := []string{}

	where = append(where, "(StartIsd=? AND StartAs=?)")
	selArgs = append(selArgs, originAS.ISD(), originAS.AS())
	updArgs = append(updArgs, originAS.ISD(), originAS.AS())
	//StartIntfGroup, AlgorithmHash
	where = append(where, "(AlgorithmHash=? AND AlgorithmID=?)")
	selArgs = append(selArgs, algHash, algId)
	updArgs = append(updArgs, algHash, algId)
	if !ignoreIntfGroup {
		where = append(where, "(StartIntfGroup=?)")
		selArgs = append(selArgs, originIntfGroup)
		updArgs = append(updArgs, originIntfGroup)
	}
	if len(where) > 0 {
		selQuery += "\n" + fmt.Sprintf("WHERE %s", strings.Join(where, " AND\n"))
		updQuery += "\n" + fmt.Sprintf("WHERE %s", strings.Join(where, " AND\n"))
	}

	if maximum > 0 {
		selQuery += "\n LIMIT ?"
		updQuery += "\n LIMIT ?"
		selArgs = append(selArgs, maximum)
		updArgs = append(updArgs, maximum)
	}
	selQuery += "\n" + "ORDER BY LastUpdated DESC"
	return selQuery, selArgs, updQuery, updArgs
}

// getBeaconMeta gets the metadata for existing beacons.
func (e *executor) getBeaconMeta(ctx context.Context, fullID []byte) (*beaconMeta, error) {
	var rowID, infoTime, lastUpdated int64
	query := "SELECT RowID, InfoTime, LastUpdated FROM Beacons WHERE FullID=?"
	err := e.db.QueryRowContext(ctx, query, fullID).Scan(&rowID, &infoTime, &lastUpdated)
	// New beacons are not in the table.
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, db.NewReadError("Failed to lookup beacon", err)
	}
	meta := &beaconMeta{
		RowID:       rowID,
		InfoTime:    time.Unix(infoTime, 0),
		LastUpdated: time.Unix(0, lastUpdated),
	}
	return meta, nil
}

// updateExistingBeacon updates the changeable data for an existing beacon.
func (e *executor) updateExistingBeacon(
	ctx context.Context,
	b beacon.Beacon,
	usage beacon.Usage,
	rowID int64,
	now time.Time,
) error {

	fullID := b.Segment.IRECID()
	packedSeg, err := beacon.PackBeaconIREC(b.Segment)
	if err != nil {
		return err
	}
	infoTime := b.Segment.Info.Timestamp.Unix()
	lastUpdated := now.UnixNano()
	expTime := b.Segment.MaxExpiry().Unix()

	// Don't have to update algorithmId, hash etc as different ones would have caused
	// different FullID and therefore RowID.
	inst := `UPDATE Beacons SET FullID=?, InIntfID=?, HopsLength=?, InfoTime=?,
			ExpirationTime=?, LastUpdated=?, Usage=?, Beacon=?, Marker = 0
			WHERE RowID=?`
	_, err = e.db.ExecContext(ctx, inst, fullID, b.InIfId, len(b.Segment.ASEntries), infoTime,
		expTime, lastUpdated, usage, packedSeg, rowID)
	if err != nil {
		return db.NewWriteError("update segment", err)
	}
	return nil
}

func insertNewBeacon(
	ctx context.Context,
	tx *sql.Tx,
	b beacon.Beacon,
	usage beacon.Usage,
	now time.Time,
) error {

	segID := b.Segment.ID()
	fullID := b.Segment.IRECID()

	packed, err := beacon.PackBeaconIREC(b.Segment)

	if err != nil {
		return db.NewInputDataError("pack segment", err)
	}
	start := b.Segment.FirstIA()
	infoTime := b.Segment.Info.Timestamp.Unix()
	lastUpdated := now.UnixNano()
	expTime := b.Segment.MaxExpiry().Unix()
	intfGroup := uint16(0)
	algorithmHash := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9} // Fallback RAC.
	algorithmId := uint32(0)
	if b.Segment.ASEntries[0].Extensions.Irec != nil {
		intfGroup = b.Segment.ASEntries[0].Extensions.Irec.InterfaceGroup
		algorithmHash = b.Segment.ASEntries[0].Extensions.Irec.AlgorithmHash
		algorithmId = b.Segment.ASEntries[0].Extensions.Irec.AlgorithmId
	}

	// Insert beacon.
	inst := `
	INSERT INTO Beacons (SegID, FullID, StartIsd, StartAs, StartIntfGroup, AlgorithmHash, AlgorithmId, InIntfID, HopsLength, InfoTime,
		ExpirationTime, LastUpdated, Usage, Beacon, Marker)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0)
	`

	// TODO(jvb): just move Irec data into the PCB, and add into the ASentry an extension
	// with a single bit for algSupported. (although it is no longer protected with a signature)

	_, err = tx.ExecContext(ctx, inst, segID, fullID, start.ISD(), start.AS(),
		intfGroup,
		algorithmHash,
		algorithmId, b.InIfId,
		len(b.Segment.ASEntries), infoTime, expTime, lastUpdated, usage, packed)
	if err != nil {
		return db.NewWriteError("insert beacon", err)
	}
	return nil
}

func (e *executor) DeleteExpiredBeacons(ctx context.Context, now time.Time) (int, error) {
	return e.deleteInTx(ctx, func(tx *sql.Tx) (sql.Result, error) {
		delStmt := `DELETE FROM Beacons WHERE ExpirationTime < ?`
		return tx.ExecContext(ctx, delStmt, now.Unix())
	})
}

func (e *executor) deleteInTx(
	ctx context.Context,
	delFunc func(tx *sql.Tx) (sql.Result, error),
) (int, error) {

	e.Lock()
	defer e.Unlock()
	return db.DeleteInTx(ctx, e.db, delFunc)
}
