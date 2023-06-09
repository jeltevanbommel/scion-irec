// Copyright 2021 Anapaya Systems
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

package routemgr

import "github.com/scionproto/scion/gateway/control"

// Dummy is a route management backend that does nothing. It is useful to use instead of nil
// pointer. It saves the user the pain of checking whether interface is nil which happens to
// be non-trivial in Go.
type Dummy struct {
}

func (d *Dummy) NewPublisher() control.Publisher {
	return d
}

func (d *Dummy) NewConsumer() control.Consumer {
	return d
}

func (d *Dummy) AddRoute(route control.Route) {}

func (d *Dummy) DeleteRoute(route control.Route) {}

func (d *Dummy) Updates() <-chan control.RouteUpdate {
	return nil
}

func (d *Dummy) Close() {}
