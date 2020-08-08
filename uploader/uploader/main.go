// +build ignore

// This is used to compile a binary for the windows special build
// (just the uploader and not the whole psy toolkit)

package main

import "github.com/psyomn/psy/uploader"

func main() {
	uploader.Run(nil)
}
