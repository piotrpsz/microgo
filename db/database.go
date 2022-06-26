package db

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"mend/db/memo"
	"mend/db/tp"
)

const (
	_ tp.DatabaseEngine = iota
	// InMemory database in memory
	InMemory
)

// Engine all databases implementations must conform to Engine
type Engine interface {
	GetAll(tp.TableName) ([]tp.Row, error)
	GetWithID(tp.TableName, tp.ID) (tp.Row, error)
	Add(tp.TableName, tp.Row) (int64, error)
	Update(tp.TableName, tp.Row) error
	Delete(tp.TableName, tp.ID) error
	Count(tp.TableName) int64
	Reset()
}

var (
	instance Engine
	once     sync.Once
	dbEngine tp.DatabaseEngine
)

// Use user can select db engine (if implemented).
// This function must be call befor first call to db.
func Use(engine tp.DatabaseEngine) {
	switch engine {
	case InMemory:
		dbEngine = engine
	default:
		log.Fatal("the db engine not implemented")
	}
}

// Db returns used current database conform to Engine.
// This implements database as signleton, is created
// only on first call.
func Db() Engine {
	once.Do(func() {
		switch dbEngine {
		case InMemory:
			instance = memo.NewInMemory()
		default:
			// no database - nothing to do
			log.Fatal("db engine not selected")
		}
	})
	return instance
}
