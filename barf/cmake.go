/*
Package barf contains code relevant to barfing.

Copyright 2019 Simon Symeonidis (psyomn)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package barf

import (
	"fmt"

	"github.com/psyomn/psy/common"
)

func cmake(_ common.RunParams) common.RunReturn {
	fmt.Println("cmake code generation goes here")
	return nil
}
