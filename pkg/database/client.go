package database

type DbClient interface {
	InsertOne() (bool, error)
	InsertMany() (bool, error)
	readOne() (bool, error)
	readMany() (bool, error)
}
