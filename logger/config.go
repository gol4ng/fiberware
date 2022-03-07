package logger

import "github.com/gol4ng/logger"

type Config struct {
	InjectOnContext bool
	LevelFunc       func(statusCode int) logger.Level
}
