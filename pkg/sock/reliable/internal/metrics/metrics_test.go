// Copyright 2019 ETH Zurich
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

package metrics_test

import (
	"testing"

	"github.com/scionproto/scion/pkg/private/prom/promtest"
	"github.com/scionproto/scion/pkg/sock/reliable/internal/metrics"
)

func TestLabels(t *testing.T) {
	promtest.CheckLabelsStruct(t, metrics.DialLabels{})
	promtest.CheckLabelsStruct(t, metrics.RegisterLabels{})
	promtest.CheckLabelsStruct(t, metrics.IOLabels{})
}
