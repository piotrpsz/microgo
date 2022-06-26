package db

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"mend/db/memo"
	"mend/db/tp"
)

const (
	_ tp.DatabasEngine = iota
	InMemory
)

type Engine interface {
	GetAll(tp.TableName) ([]tp.Row, error)
	GetWithID(tp.TableName, tp.ID) (tp.Row, error)
	Add(tp.TableName, tp.Row) (int64, error)
	Update(tp.TableName, tp.Row) error
	Delete(tp.TableName, tp.ID) error
	Count(tp.TableName) int64
}

var (
	instance Engine
	once     sync.Once
	dbEngine tp.DatabasEngine
)

// Use user can select db engine (if implemented).
// This function must be call befor first call to db.
func Use(engine tp.DatabasEngine) {
	switch engine {
	case InMemory:
		dbEngine = engine
	default:
		log.Fatal("the db engine not implemented")
	}
}

func Db() Engine {
	once.Do(func() {
		switch dbEngine {
		case InMemory:
			instance = memo.NewInMemory()
		default:
			log.Fatal("db engine not selected")
		}
	})
	return instance
}
