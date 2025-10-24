package log

import "go.uber.org/zap/zapcore"

type Options struct {
	Level         zapcore.Level
	LogDir        string
	LineNum       bool
	CallerSkipNum int
}

func NewOptions() *Options {

	return &Options{}
}
