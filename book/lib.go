/*
Package book is a simple convenient fetcher of books.

https://www.goodreads.com/shelf/list.xml?key=XXXXXXXXXXXXXXXXXXXXX
*/
package book

import (
	"github.com/psyomn/psy/common"
)

// Run will check the book list of the user
func Run(args common.RunParams) common.RunReturn {
	err := loadConfig()
	if err != nil {
		return err
	}

	err = authenticate()
	if err != nil {
		return err
	}

	return nil
}
