// Copyright 2013 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hlog

import (
	"flag"
)

func main() {

	// init output dir
	flag.Parse()

	// API:Print
	Print("error", "the error code/message: ", 400, "/", "bad request")

	// API::Printf
	Printf("error", "the error code/message: %d/%s", 400, "bad request")
}
