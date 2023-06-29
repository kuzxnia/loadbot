package database

type DbClient interface {
	InsertOne() (bool, error)
	InsertMany() (bool, error)
	InsertOneOrMany() (bool, error)
	ReadOne() (bool, error)
	ReadMany() (bool, error)

  GetBatchSize() uint64
}
