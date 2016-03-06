package orm

import (
	"strconv"

	"gopkg.in/pg.v4/types"
)

type dber interface {
	Exec(q interface{}, params ...interface{}) (types.Result, error)
	ExecOne(q interface{}, params ...interface{}) (types.Result, error)
	Query(coll, query interface{}, params ...interface{}) (types.Result, error)
	QueryOne(model, query interface{}, params ...interface{}) (types.Result, error)
}

type selectQuery struct {
	*Query
}

func (sel selectQuery) AppendQuery(b []byte, params ...interface{}) ([]byte, error) {
	b = append(b, "SELECT "...)
	if sel.columns == nil {
		b = types.AppendField(b, sel.model.Table.Name, true)
		b = append(b, ".*"...)
	} else {
		b = appendValue(b, ", ", sel.columns...)
	}

	b = append(b, " FROM "...)
	b = appendField(b, sel.tables...)

	b = appendString(b, "", sel.joins...)

	if sel.wheres != nil {
		b = append(b, " WHERE "...)
		b = appendString(b, " AND ", sel.wheres...)
	}

	if sel.orders != nil {
		b = append(b, " ORDER BY "...)
		b = appendString(b, ", ", sel.orders...)
	}

	if sel.limit != 0 {
		b = append(b, " LIMIT "...)
		b = strconv.AppendInt(b, int64(sel.limit), 10)
	}

	if sel.offset != 0 {
		b = append(b, " OFFSET "...)
		b = strconv.AppendInt(b, int64(sel.offset), 10)
	}

	return b, nil
}

func appendField(b []byte, ss ...string) []byte {
	for i, field := range ss {
		b = types.AppendField(b, field, true)
		if i != len(ss)-1 {
			b = append(b, ", "...)
		}
	}
	return b
}

func appendString(b []byte, sep string, ss ...string) []byte {
	for i, s := range ss {
		b = append(b, s...)
		if i != len(ss)-1 {
			b = append(b, sep...)
		}
	}
	return b
}

func appendValue(b []byte, sep string, vv ...types.ValueAppender) []byte {
	for i, v := range vv {
		var err error
		b, err = v.AppendValue(b, true)
		if err != nil {
			panic(err)
		}
		if i != len(vv)-1 {
			b = append(b, sep...)
		}
	}
	return b
}
