package filter

import (
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type filterStruct struct {
	Name  string
	Expr  string
	Value interface{}
}

func Parse(filter interface{}) ([]filterStruct, error) {
	if reflect.TypeOf(filter).Kind() != reflect.Ptr {
		return nil, errors.New("filter must be a pointer")
	}
	var f []filterStruct
	t := reflect.TypeOf(filter).Elem()
	v := reflect.ValueOf(filter).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if _, ok := t.Field(i).Tag.Lookup("filter"); !ok {
			continue
		}
		if !field.IsValid() {
			continue
		}

		if reflect.Int <= field.Kind() && field.Kind() <= reflect.Int64 && field.Int() == 0 {
			continue
		}
		if reflect.Uint <= field.Kind() && field.Kind() <= reflect.Uint64 && field.Uint() == 0 {
			continue
		}
		if reflect.Float32 <= field.Kind() && field.Kind() <= reflect.Float64 && field.Float() == 0 {
			continue
		}
		if field.Kind() == reflect.String && field.String() == "" {
			continue
		}
		if field.Kind() == reflect.Pointer {
			if field.IsNil() {
				continue
			}
			field = reflect.Indirect(field)
		}

		fs := filterStruct{
			Value: field.Interface(),
		}

		tags := strings.Split(t.Field(i).Tag.Get("filter"), ";")
		if len(tags) > 2 {
			return nil, errors.New("filter tag error")
		} else if len(tags) == 2 {
			fs.Name = strings.Split(tags[0], ":")[1]
			if len(tags[1]) > 0 {
				fs.Expr = strings.Split(tags[1], ":")[1]
			}
		} else {
			fs.Name = strings.Split(tags[0], ":")[1]
		}
		f = append(f, fs)

	}
	return f, nil
}

func Query(db *gorm.DB, filter interface{}) (*gorm.DB, error) {
	// log.Printf("%#v\n", filter)
	f, err := Parse(filter)
	if err != nil {
		return nil, err
	}
	exprs := make([]clause.Expression, 0, len(f))
	for _, v := range f {
		switch v.Expr {
		case "exact":
			exprs = append(exprs, clause.Eq{Column: v.Name, Value: v.Value})
		case "iexact":
			exprs = append(exprs, IExact{Column: v.Name, Value: v.Value})
		case "contains":
			exprs = append(exprs, Contains{Column: v.Name, Value: v.Value})
		case "icontains":
			exprs = append(exprs, IContains{Column: v.Name, Value: v.Value})
		case "in":
			db = db.Clauses(clause.IN{Column: v.Name, Values: v.Value.([]interface{})})
		case "gt":
			db = db.Clauses(clause.Gt{Column: v.Name, Value: v.Value})
		case "gte":
			db = db.Clauses(clause.Gte{Column: v.Name, Value: v.Value})
		case "lt":
			db = db.Clauses(clause.Lt{Column: v.Name, Value: v.Value})
		case "lte":
			db = db.Clauses(clause.Lte{Column: v.Name, Value: v.Value})
		case "startswith":
			db = db.Clauses(StartsWith{Column: v.Name, Value: v.Value})
		case "istartswith":
			db = db.Clauses(IStartsWith{Column: v.Name, Value: v.Value})
		case "endswith":
			db = db.Clauses(EndsWith{Column: v.Name, Value: v.Value})
		case "iendswith":
			db = db.Clauses(IEndsWith{Column: v.Name, Value: v.Value})
		default:
			db = db.Clauses(clause.Eq{Column: v.Name, Value: v.Value})
		}
	}
	db.Clauses(exprs...)
	return db, nil
}
