package spoukfw

import (
	"log"
	"io"
	"os"
	"fmt"
)

type (
	SpoukLogger struct {
		logger  *log.Logger
		Logging io.Writer
	}
)

const (
	loggerPrefix  = "[spoukftw][spouklogger][%s]"
)
//# TODO[spouk-26.08.16] добавить цветовую палитру по типу сообщений
//# TODO[spouk-26.08.16] добавить поддержку множественности логгирования с разными каналами трансляции

func NewSpoukLogger(subprefix string, logging io.Writer) *SpoukLogger {
	sl := &SpoukLogger{Logging:logging}

	if sl.Logging == nil {
		sl.Logging = os.Stdout
	}
	sl.logger = log.New(sl.Logging, fmt.Sprintf(loggerPrefix, subprefix), 0)
	return sl
}
func (s *SpoukLogger) Printf(format string, v ...interface{}) {
	s.logger.Printf(format, v)
}
func (s *SpoukLogger) Info(message string) {
	s.logger.Printf(message)
}
func (s *SpoukLogger) Error(message string) {
	s.logger.Printf(message)
}
func (s *SpoukLogger) Warning(message string) {
	s.logger.Printf(message)
}