package main

import (
	"fmt"
	"sort"
	"time"
)

const defaultTimeFormat = "2006-01-02T15:04:05Z" //oddly enough this is how you specify a format

const (
	//TypeInt is printed as %d
	TypeInt = iota
	//TypeString is printed as %s
	TypeString = iota
	//TypeFloat is printed as %f
	TypeFloat = iota
	//TypeDateTime is printed as a string after parsing
	TypeDateTime = iota
	//TypeInterface is printed as %v
	TypeInterface = iota
	//TypeBool is printed as %v
	TypeBool = iota
)

//SchemaField defines a field in a table
type SchemaField struct {
	FieldName      string
	FieldType      int
	FieldSize      int
	FieldPrecision int
	FieldFormat    string
}

type lessFunc func(p1, p2 interface{}, field *SchemaField) bool

// MultiSorter implements the Sort interface, sorting the changes within.
type MultiSorter struct {
	data    [][]interface{}
	less    []lessFunc
	schema  []SchemaField
	indexes []int
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *MultiSorter) Sort(data [][]interface{}) {
	ms.data = data
	sort.Sort(ms)
}

// Len is part of sort.Interface.
func (ms *MultiSorter) Len() int {
	return len(ms.data)
}

// Swap is part of sort.Interface.
func (ms *MultiSorter) Swap(i, j int) {
	ms.data[i], ms.data[j] = ms.data[j], ms.data[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that discriminates between
// the two items (one is less than the other). Note that it can call the
// less functions twice per call. We could change the functions to return
// -1, 0, 1 and reduce the number of calls for greater efficiency: an
// exercise for the reader.
func (ms *MultiSorter) Less(i, j int) bool {
	p, q := ms.data[i], ms.data[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		index := ms.indexes[k]
		switch {
		case less(p[index], q[index], &ms.schema[index]):
			// p < q, so we have a decision.
			return true
		case less(q[index], p[index], &ms.schema[index]):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	lastIndex := ms.indexes[k]
	return ms.less[k](p[lastIndex], q[lastIndex], &ms.schema[lastIndex])
}

//TableSorter a multisorter for a table
func TableSorter(schema []SchemaField) *MultiSorter {
	return &MultiSorter{
		schema: schema,
	}
}

//OrderBy specifies the order
func (ms *MultiSorter) OrderBy(fieldNames ...string) *MultiSorter {

	ms.less = make([]lessFunc, len(fieldNames))
	ms.indexes = make([]int, len(fieldNames))

	for k, fn := range fieldNames {
		var field *SchemaField
		for i, f := range ms.schema {

			if f.FieldName == fn {
				field = &f
				ms.indexes[k] = i
				break
			}
		}

		if field == nil {
			fmt.Printf("could not find field with name %s\n", fn)
			return nil
		}

		switch field.FieldType {
		case TypeInt:
			ms.less[k] = func(a, b interface{}, field *SchemaField) bool {
				return a.(int) < b.(int)
			}
		case TypeString:
			ms.less[k] = func(a, b interface{}, field *SchemaField) bool {
				return a.(string) < b.(string)
			}
		case TypeFloat:
			ms.less[k] = func(a, b interface{}, field *SchemaField) bool {
				return a.(float64) < b.(float64)
			}
		case TypeDateTime:
			ms.less[k] = func(a, b interface{}, field *SchemaField) bool {

				layout := defaultTimeFormat

				if field.FieldFormat != "" {
					layout = field.FieldFormat
				}

				ta, err := time.Parse(layout, a.(string))
				if err != nil {
					fmt.Printf("could not convert string %s to date time", a.(string))
					return false
				}

				tb, err := time.Parse(layout, b.(string))
				if err != nil {
					fmt.Printf("could not convert string %s to date time", b.(string))
					return false
				}

				return ta.Before(tb)
			}
		case TypeBool:
			ms.less[k] = func(a, b interface{}, field *SchemaField) bool {
				return a.(bool) != b.(bool)
			}
		default:
			fmt.Printf("could not find type %d", field.FieldType)
		}
	}

	return ms
}
