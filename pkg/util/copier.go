package util

import "github.com/gotidy/copy"

var _copier = copy.New()

func Copy[T, R any](src T, dest R) {
	_copier.Copy(src, dest)
}
