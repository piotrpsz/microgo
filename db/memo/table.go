package memo

import (
	"container/list"
	"sync"

	"mend/db/tp"
)

// MemoryTable data container (table).
// Component of InMemory
type MemoryTable struct {
	sync.Mutex
	data  *list.List
	index map[tp.ID]*list.Element
	maxID tp.ID
}

// NewTable creates new instance of memory-table.
func NewTable() *MemoryTable {
	return &MemoryTable{
		data:  list.New(),
		index: make(map[tp.ID]*list.Element),
	}
}

func (t *MemoryTable) reset() {
	t.Lock()
	defer t.Unlock()

	t.data = list.New()
	t.index = make(map[tp.ID]*list.Element)
}

func (t *MemoryTable) getAll() []tp.Row {
	t.Lock()
	defer t.Unlock()

	data := make([]tp.Row, 0, t.data.Len())
	for element := t.data.Front(); element != nil; element = element.Next() {
		data = append(data, element.Value.(tp.Row))
	}
	return data
}

func (t *MemoryTable) getWithID(id tp.ID) (tp.Row, error) {
	t.Lock()
	defer t.Unlock()

	element, ok := t.index[id]
	if !ok {
		return nil, tp.ErrRowNotFound
	}
	return element.Value.(tp.Row), nil
}

func (t *MemoryTable) add(data tp.Row) tp.ID {
	idx := tp.ID(len(t.index) + 1)
	data["id"] = idx

	t.Lock()
	defer t.Unlock()

	entry := t.data.PushBack(data)
	t.index[idx] = entry

	return idx
}

func (t *MemoryTable) update(data tp.Row) error {
	if id, ok := data["id"]; ok {
		if idv, ok := id.(tp.ID); ok {
			t.Lock()
			defer t.Unlock()

			if element, ok := t.index[idv]; ok {
				element.Value = data
				return nil
			}
			return tp.ErrRowNotFound
		}
	}
	return tp.ErrInvalidRowData
}

func (t *MemoryTable) delete(id tp.ID) error {
	t.Lock()
	defer t.Unlock()

	element, ok := t.index[id]
	if !ok {
		return tp.ErrRowNotFound
	}

	t.data.Remove(element)
	delete(t.index, id)

	return nil
}

func (t *MemoryTable) count() int64 {
	t.Lock()
	defer t.Unlock()

	return int64(len(t.index))
}
