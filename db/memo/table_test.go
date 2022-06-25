package memo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mend/db/tp"
)

func TestTable_add(t *testing.T) {
	var data = []struct {
		row tp.Row
		id  tp.ID
	}{
		{
			tp.Row{"name": "a", "age": 100},
			tp.ID(1),
		},
		{
			tp.Row{"id": 0, "name": "b", "age": 101},
			tp.ID(2),
		},
		{
			tp.Row{"name": "a", "age": 100},
			tp.ID(3),
		},
		{
			tp.Row{"id": 0, "name": "c", "age": 102},
			tp.ID(4),
		},
		{
			tp.Row{"id": 100, "name": "b", "age": 101},
			tp.ID(5),
		},
	}

	table := NewTable()
	for _, tt := range data {
		id := table.add(tt.row)
		assert.Equal(t, tt.id, id)
	}
}

func TestTable_getWithID(t *testing.T) {
	var data = []tp.Row{
		{"name": "a", "age": 100},
		{"name": "b", "age": 101},
		{"name": "b", "age": 102},
		{"name": "b", "age": 103},
		{"name": "b", "age": 104},
		{"name": "b", "age": 105},
	}

	table := NewTable()
	for _, tt := range data {
		id := table.add(tt)
		tt["id"] = id
	}

	// positive tests
	for _, tt := range data {
		id := tt["id"].(tp.ID)
		result, err := table.getWithID(id)
		assert.NoError(t, err)
		assert.Equal(t, result, tt)
	}

	// negative tests
	for i := 100; i < 110; i++ {
		result, err := table.getWithID(tp.ID(i))
		assert.Error(t, err, tp.ErrRowNotFound)
		assert.Nil(t, result)
	}
}

func TestTable_getAll(t *testing.T) {
	var data = [][]tp.Row{
		{
			{"id": tp.ID(1), "name": "a", "age": 100},
		},
		{
			{"id": tp.ID(1), "name": "a", "age": 100},
			{"id": tp.ID(2), "name": "b", "age": 101},
		},
		{
			{"id": tp.ID(1), "name": "a", "age": 100},
			{"id": tp.ID(2), "name": "b", "age": 101},
			{"id": tp.ID(3), "name": "b", "age": 102},
		},
		{
			{"id": tp.ID(1), "name": "a", "age": 100},
			{"id": tp.ID(2), "name": "b", "age": 101},
			{"id": tp.ID(3), "name": "b", "age": 102},
			{"id": tp.ID(4), "name": "b", "age": 103},
		},
		{
			{"id": tp.ID(1), "name": "a", "age": 100},
			{"id": tp.ID(2), "name": "b", "age": 101},
			{"id": tp.ID(3), "name": "b", "age": 102},
			{"id": tp.ID(4), "name": "b", "age": 103},
			{"id": tp.ID(5), "name": "b", "age": 104},
		},
		{
			{"id": tp.ID(1), "name": "a", "age": 100},
			{"id": tp.ID(2), "name": "b", "age": 101},
			{"id": tp.ID(3), "name": "b", "age": 102},
			{"id": tp.ID(4), "name": "b", "age": 103},
			{"id": tp.ID(5), "name": "b", "age": 104},
			{"id": tp.ID(6), "name": "b", "age": 105},
		},
	}

	for _, tt := range data {
		table := NewTable()
		for _, row := range tt {
			_ = table.add(row)
		}

		result := table.getAll()
		assert.Equal(t, tt, result)
	}
}

func TestTable_update(t *testing.T) {
	var data = []struct {
		src tp.Row
		dst []tp.Row
		err error
	}{
		{
			tp.Row{"id": tp.ID(1), "name": "a", "age": 100},
			[]tp.Row{
				{"id": tp.ID(1), "name": "b", "age": 101},
				{"id": tp.ID(1), "name": "c", "age": 102},
				{"id": tp.ID(1), "name": "d", "age": 103},
			},
			nil,
		},
	}

	table := NewTable()
	for _, tt := range data {
		tt.src["id"] = table.add(tt.src)
	}

	for _, tt := range data {
		id := tt.src["id"].(tp.ID)
		for _, item := range tt.dst {
			err := table.update(item)
			assert.Equal(t, tt.err, err)

			row, err := table.getWithID(id)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, item, row)
		}
	}
}
