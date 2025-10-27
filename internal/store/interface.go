package store

type DataBase interface {
	Open()
	Close()
	Ping()
	GetDB() any
}
