package repository

type Repository interface {
	InsertRecord(any) error
	SelectRecord(any) error
	DeleteRecord(any) error
	UpdateRecord(any) error
}
