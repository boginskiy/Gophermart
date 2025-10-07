package repository

type Repository interface {
	CheckUnicRecord(record any) bool
	InsertRecord(any) error
	// SelectRecord(any) error
	// DeleteRecord(any) error
	// UpdateRecord(any) error
}
