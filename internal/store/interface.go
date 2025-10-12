package store

type Dber interface {
	Open()
	Clouse()
	Ping()
	GetDB() any
}
