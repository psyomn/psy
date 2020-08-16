/*
Package common contains things that should be exported anywhere within
psy.

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
package common

import (
	"os"
	"path"
)

// RunParams is the type signature that all run commands should respect
type RunParams = []string

// RunReturn is the return type that all run commands should respect
type RunReturn = error

// ConfigDir prefers $HOME/.config, regardles of XDG stuff (for now)
func ConfigDir() string {
	homeDir := os.Getenv("HOME")

	if homeDir == "" {
		panic("need home to run")
	}

	return path.Join(homeDir, ".config", "psy")
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer fd.Close()
	return err != nil
}
