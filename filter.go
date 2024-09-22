package pglogicalstream

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ChangeFilter struct {
	tablesWhiteList map[string]bool
	schemaWhiteList string
}

type Filtered func(change Wal2JsonChanges)

func NewChangeFilter(tableSchemas []string, schema string) ChangeFilter {
	tablesMap := map[string]bool{}
	for _, table := range tableSchemas {
		tablesMap[table] = true
	}

	return ChangeFilter{
		tablesWhiteList: tablesMap,
		schemaWhiteList: schema,
	}
}

func (c ChangeFilter) FilterChange(lsn string, change []byte, OnFiltered Filtered) {
	var changes WallMessage
	if err := json.NewDecoder(bytes.NewReader(change)).Decode(&changes); err != nil {
		panic(fmt.Errorf("cant parse change from database to filter it %v", err))
	}

	if len(changes.Change) == 0 {
		return
	}

	for _, ch := range changes.Change {
		var filteredChanges = Wal2JsonChanges{
			Lsn:     &lsn,
			Changes: []Wal2JsonChange{},
		}
		if ch.Schema != c.schemaWhiteList {
			continue
		}

		var (
			tableExist bool
		)
		fmt.Println(ch.Table, c.tablesWhiteList)
		if _, tableExist = c.tablesWhiteList[ch.Table]; !tableExist {
			continue
		}

		if ch.Kind == "delete" {
			for i, changedValue := range ch.Oldkeys.Keyvalues {
				ch.Columnvalues[i] = changedValue
			}
		}

		filteredChanges.Changes = append(filteredChanges.Changes, Wal2JsonChange{
			Kind:         ch.Kind,
			Schema:       ch.Schema,
			Table:        ch.Table,
			ColumnNames:  ch.Columnnames,
			ColumnTypes:  ch.Columntypes,
			ColumnValues: ch.Columnvalues,
		})

		OnFiltered(filteredChanges)
	}
}
