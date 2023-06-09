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
	"flag"
	"fmt"
	"os"
)

const (
	DefaultTopoFile = "topology/default.topo"
	DefaultGenFile  = "pkg/private/xtest/graph/default_gen.go"
)

var (
	topoFile  = flag.String("topoFile", DefaultTopoFile, "")
	graphFile = flag.String("graphFile", DefaultGenFile, "")
)

func main() {
	err := WriteGraphToFile(*topoFile, *graphFile)
	if err != nil {
		fmt.Printf("Failed to write the graph, err: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("Successfully written graph to %s\n", *graphFile)
	}
}
