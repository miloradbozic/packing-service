package database

// PackSizeRepositoryInterface defines the interface for pack size repository operations
type PackSizeRepositoryInterface interface {
	GetAll() ([]PackSize, error)
	GetByID(id int) (*PackSize, error)
	Create(size int) (*PackSize, error)
	Update(id int, size int) (*PackSize, error)
	Delete(id int) error
}
