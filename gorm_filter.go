package filter

import (
	"errors"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type filterStruct struct {
	Name  string
	Expr  string
	Value interface{}
}

func Parse(filter interface{}) ([]filterStruct, error) {
	var f []filterStruct
	t := reflect.TypeOf(filter).Elem()
	v := reflect.ValueOf(filter).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		//log.Printf("%s field kind: %+v\n", t.Field(i).Name, field.Kind())
		if !v.Field(i).IsValid() { // not valid
			continue
		}

		if field.Kind() == reflect.Pointer {
			if field.IsNil() { // nil value
				continue
			}
			field = reflect.Indirect(field)
		}

		if field.Len() == 0 {
			continue
		}
		tags := strings.Split(t.Field(i).Tag.Get("filter"), ";")
		if len(tags) != 2 {
			return nil, errors.New("filter tag error")
		}
		f = append(f, filterStruct{
			Name:  strings.Split(tags[0], ":")[1],
			Expr:  strings.Split(tags[1], ":")[1],
			Value: field.Interface(),
		})
	}
	return f, nil
}

func Query(db *gorm.DB, filter interface{}) (*gorm.DB, error) {
	f, err := Parse(filter)
	if err != nil {
		return nil, err
	}
	for _, v := range f {
		if v.Value == nil {
			continue
		}
		switch v.Expr {
		case "exact":
			// db = db.Clauses(clause.Eq{Column: v.Name, Value: v.Value})
			db = db.Where(v.Name+" = ?", v.Value)
		case "iexact":
			db = db.Where(v.Name+" ILIKE ?", v.Value)
		case "like", "contains":
			db = db.Where(v.Name+" LIKE ?", "%"+v.Value.(string)+"%")
		case "ilike", "icontains":
			db = db.Where(v.Name+" ILIKE ?", "%"+v.Value.(string)+"%")
		case "in":
			db = db.Where(v.Name+" IN ?", v.Value)
		case "gt":
			db = db.Where(v.Name+" > ?", v.Value)
		case "gte":
			db = db.Where(v.Name+" >= ?", v.Value)
		case "lt":
			db = db.Where(v.Name+" < ?", v.Value)
		case "lte":
			db = db.Where(v.Name+" <= ?", v.Value)
		case "startswith":
			db = db.Where(v.Name+" LIKE ?", v.Value.(string)+"%")
		case "istartswith":
			db = db.Where(v.Name+" ILIKE ?", v.Value.(string)+"%")
		case "endswith":
			db = db.Where(v.Name+" LIKE ?", "%"+v.Value.(string))
		case "iendswith":
			db = db.Where(v.Name+" ILIKE ?", "%"+v.Value.(string))
		case "isnull":
			db = db.Where(v.Name + " IS NULL")
		default:
			db = db.Where(v.Name+" = ?", v.Value)
		}
	}
	return db, nil
}
