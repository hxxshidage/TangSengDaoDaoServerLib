package util

import "github.com/gotidy/copy"

var _copier = copy.New(copy.Skip())

func Copy[T, R any](src T, dest R) R {
	_copier.Copy(src, dest)
	return dest
}
