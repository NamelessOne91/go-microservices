package data

type Repository interface {
	Insert(LogEntry) error
	All() ([]*LogEntry, error)
	GetOne(id string) (*LogEntry, error)
	DropCollection() error
	Update(LogEntry) (bool, error)
}
