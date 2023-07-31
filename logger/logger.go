package logger

import (
	"log"
	"os"
)

const (
	infoL  = "[INFO] "
	errL   = "[ERROR] "
	fatalL = "[FATAL] "
)

const (
	ief = iota
	ie
	i
	ef
	e
	f
)

type logger struct {
	logger *log.Logger
	level  int
}

func (l *logger) Info(v ...any) {
	if l.level < 3 {
		l.logger.SetPrefix(infoL)
		l.logger.Println(v...)
	}
}

func (l *logger) Error(v ...any) {

	if l.level != 2 && l.level != 5  {
	l.logger.SetPrefix(errL)
	l.logger.Println(v...)
	}
}

func (l *logger) Fatal(v ...any) {
	l.logger.SetPrefix(fatalL)
	l.logger.Fatalln(v...)
}

func NewLogger() *logger {
	return &logger{
		log.New(os.Stdout, "", 19),
		ief,
	}
}
