package logg

type Logger interface {
	RaiseError(msg string, err error)
	RaiseFatal(msg string, err error)
	RaisePanic(msg string, err error)
	RaiseInfo(msg string)
	Close()
}
