package repository

type Repository interface {
	CheckUnicRecord(record any) bool
	InsertRecord(record any) error
	SelectRecord(any) (record any, err error)
	// DeleteRecord(any) error
	// UpdateRecord(any) error
}
