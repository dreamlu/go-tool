package excel

import (
	"github.com/dreamlu/gt/src/reflect"
	"github.com/dreamlu/gt/src/type/tmap"
	"github.com/xuri/excelize/v2"
	"io"
)

func (f *Excel[T]) AddDict(key string, value dict) *Excel[T] {
	f.dict.Set(key, value)
	return f
}

func (f *Excel[T]) Import(r io.Reader, opts ...excelize.Options) (datas []*T, err error) {

	f.File, err = excelize.OpenReader(r, opts...)
	if err != nil {
		return
	}
	defer f.Close()
	rows, err := f.GetRows(f.sheet)
	if err != nil {
		return nil, err
	}

	var (
		title = tmap.NewTMap[string, int]()
		max   = len(rows[0])
	)
	for k, colCell := range rows[0] {
		title.Set(colCell, k)
	}

	for i := 1; i < len(rows); i++ {

		row := rows[i]
		for len(row) < max {
			row = append(row, "")
		}
		var data T
		for k, v := range f.ExcelMapper {
			if fc := f.dict.Get(v); fc != nil {
				var value any
				value, err = fc(v, row[title.Get(v)])
				if err != nil {
					return
				}
				reflect.Set(&data, k.Field, value)
				continue
			}
			if !title.IsExist(v) {
				continue
			}
			value := string2any(k.Type, row[title.Get(v)])
			reflect.Set(&data, k.Field, value)
		}
		datas = append(datas, &data)
	}

	// after import
	var data T
	if reflect.IsImplements(data, new(Handle[T])) {
		err = reflect.Call(data, "ExcelHandle", datas)
		if err != nil {
			return
		}
	}

	return
}
