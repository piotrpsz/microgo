package memo

import (
	"fmt"

	"mend/db/tp"
)

// InMemory constains tables with data.
// Must conform to the tp.Engine interface.
type InMemory struct {
	tables map[tp.TableName]*memoryTable
}

// NewInMemory return instance od in-memory 'tp..
// It implementet as singleton: object is created only on first call.
func NewInMemory() *InMemory {
	return &InMemory{
		tables: map[tp.TableName]*memoryTable{
			"user": NewTable(),
		},
	}
}

func (m *InMemory) Count(tblName tp.TableName) int64 {
	if table, ok := m.tables[tblName]; ok {
		return table.count()
	}
	return -1

}

// GetAll returns all rows from table.
func (m *InMemory) GetAll(tblName tp.TableName) ([]tp.Row, error) {
	if table, ok := m.tables[tblName]; ok {
		return table.getAll(), nil
	}
	return nil, fmt.Errorf("table '%s' not found", tblName)
}

// GetWithID returns data of row with ID.
func (m *InMemory) GetWithID(tblName tp.TableName, id tp.ID) (tp.Row, error) {
	if table, ok := m.tables[tblName]; ok {
		return table.getWithID(id)
	}
	return nil, fmt.Errorf("table '%s' not found", tblName)
}

// Add adds new data row to table.
func (m *InMemory) Add(tblName tp.TableName, data tp.Row) (tp.ID, error) {
	if table, ok := m.tables[tblName]; ok {
		return table.add(data), nil
	}
	return -1, fmt.Errorf("table '%s' not found", tblName)
}

// Update updates in table data of row with ID.
func (m *InMemory) Update(tblName tp.TableName, data tp.Row) error {
	if table, ok := m.tables[tblName]; ok {
		return table.update(data)
	}
	return fmt.Errorf("table '%s' not found", tblName)

}

// Delete removes from table row with ID.
func (m *InMemory) Delete(tblName tp.TableName, id tp.ID) error {
	if table, ok := m.tables[tblName]; ok {
		return table.delete(id)
	}
	return fmt.Errorf("table '%s' not found", tblName)
}

func (m *InMemory) Reset() {
	if table, ok := m.tables["user"]; ok {
		table.reset()
	}
}
