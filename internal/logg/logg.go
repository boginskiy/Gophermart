package logg

import (
	"log"
	"os"

	"go.uber.org/zap"
)

type Logg struct {
	Sugar *zap.SugaredLogger
	file  *os.File
}

func NewLogg(logFile string) Logger {
	// Файл для логов
	f, err := CreateFile(logFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// Default config
	logger, err := TakeConfig(f.Name())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return &Logg{
		Sugar: logger.Sugar(),
		file:  f,
	}
}

func (l *Logg) Close() {
	l.file.Close() // Закрываем дискриптор файл
	l.Sugar.Sync() // Гарантия, что все оставшиеся сообщения попадут в лог
}

func (l *Logg) RaiseInfo(msg string) {
	l.Sugar.Debugf("[INFO]: %s", msg)
}

func (l *Logg) RaiseError(msg string, err error) {
	l.Sugar.Errorf("[ERROR]: %s | %s", err.Error(), msg)
}

func (l *Logg) RaiseFatal(msg string, err error) {
	l.Sugar.Errorf("[FATAL]: %s | %s", err.Error(), msg)
}

func (l *Logg) RaisePanic(msg string, err error) {
	l.Sugar.Errorf("[PANIC]: %s | %s", err.Error(), msg)
}
