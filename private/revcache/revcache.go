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

package revcache

import (
	"context"
	"fmt"
	"io"

	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/private/common"
	"github.com/scionproto/scion/pkg/private/ctrl/path_mgmt"
	"github.com/scionproto/scion/private/storage/db"
)

// Key denotes the key for the revocation cache.
type Key struct {
	IA   addr.IA
	IfId common.IFIDType
}

// NewKey creates a new key for the revocation cache.
func NewKey(ia addr.IA, ifId common.IFIDType) *Key {
	return &Key{
		IA:   ia,
		IfId: ifId,
	}
}

func (k Key) String() string {
	return fmt.Sprintf("%s#%s", k.IA, k.IfId)
}

// KeySet is a set of keys.
type KeySet map[Key]struct{}

// SingleKey is a convenience function to return a KeySet with a single key.
func SingleKey(ia addr.IA, ifId common.IFIDType) KeySet {
	return KeySet{*NewKey(ia, ifId): {}}
}

// RevOrErr is either a revocation or an error.
type RevOrErr struct {
	Rev *path_mgmt.RevInfo
	Err error
}

// ResultChan is a channel of results.
type ResultChan <-chan RevOrErr

// RevCache is a cache for revocations. Revcache implementations must be safe for concurrent usage.
type RevCache interface {
	// Get items with the given keys from the cache. Returns all present requested items that are
	// not expired or an error if the query failed.
	Get(ctx context.Context, keys KeySet) (Revocations, error)
	// GetAll returns a channel that will provide all items in the revocation cache. If the cache
	// can't prepare the result channel a nil channel and the error are returned. If the querying
	// succeeded the channel will contain the revocations in the cache, or an error if the stored
	// data could not be parsed. Note that implementations can spawn a goroutine to fill the
	// channel, which means the channel must be fully drained to guarantee the destruction of the
	// goroutine.
	GetAll(ctx context.Context) (ResultChan, error)
	// Insert inserts or updates the given revocation into the cache.
	// Returns whether an insert was performed.
	Insert(ctx context.Context, rev *path_mgmt.RevInfo) (bool, error)
	// DeleteExpired deletes expired entries from the cache.
	// Users of the revcache should make sure to periodically call this method to prevent an
	// ever growing cache.
	// Returns the amount of deleted entries.
	DeleteExpired(ctx context.Context) (int64, error)
	db.LimitSetter
	io.Closer
}

// Revocations is the map of revocations.
type Revocations map[Key]*path_mgmt.RevInfo

// RevocationToMap converts a slice of revocations to a revocation map. If extracting the info field
// fails from a revocation, nil and the error is returned.
func RevocationToMap(revs []*path_mgmt.RevInfo) (Revocations, error) {
	res := make(Revocations)
	for _, rev := range revs {
		res[Key{IA: rev.IA(), IfId: rev.IfID}] = rev
	}
	return res, nil
}

// Keys returns the set of keys in this revocation map.
func (r Revocations) Keys() KeySet {
	keys := make(KeySet, len(r))
	for key := range r {
		keys[key] = struct{}{}
	}
	return keys
}

// ToSlice extracts the values from this revocation map.
func (r Revocations) ToSlice() []*path_mgmt.RevInfo {
	res := make([]*path_mgmt.RevInfo, 0, len(r))
	for _, rev := range r {
		res = append(res, rev)
	}
	return res
}

// FilterNew drops all revocations in r that are already in the revCache.
func (r Revocations) FilterNew(ctx context.Context, revCache RevCache) error {
	inCache, err := revCache.Get(ctx, r.Keys())
	if err != nil {
		return err
	}
	for key, rev := range r {
		existingRev, exists := inCache[key]
		if exists {
			if !newerInfo(existingRev, rev) {
				delete(r, key)
			}
		}
	}
	return nil
}
