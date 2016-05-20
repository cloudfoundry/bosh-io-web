package index

import (
	"database/sql"
	"encoding/json"
	"reflect"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

type DBIndex struct {
	adapter DBAdapter

	logTag string
	logger boshlog.Logger
}

type dbIndexLockedRecord struct {
	rawKeyBytes []byte
	notFound    bool
	tx          *sql.Tx
	dbIndex     DBIndex

	logTag string
	logger boshlog.Logger
}

func NewDBIndex(adapter DBAdapter, logger boshlog.Logger) DBIndex {
	return DBIndex{adapter: adapter, logTag: "DBIndex", logger: logger}
}

func (ri DBIndex) List(values interface{}) error {
	rows, err := ri.adapter.List()
	if err != nil {
		return err
	}

	var rawValues []json.RawMessage

	for rows.Next() {
		var valueBytes []byte

		err = rows.Scan(&valueBytes)
		if err != nil {
			return err
		}

		rawValues = append(rawValues, valueBytes)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	// todo avoid serializing already collected entries
	rawValuesBytes, err := json.Marshal(rawValues)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawValuesBytes, values)
	if err != nil {
		return err
	}

	return nil
}

func (ri DBIndex) ListKeys(keys interface{}) error {
	rows, err := ri.adapter.ListKeys()
	if err != nil {
		return err
	}

	keysElem := reflect.ValueOf(keys).Elem()

	for rows.Next() {
		var rawKeyBytes []byte

		err = rows.Scan(&rawKeyBytes)
		if err != nil {
			return err
		}

		var rawKey map[string]interface{}

		err = json.Unmarshal(rawKeyBytes, &rawKey)
		if err != nil {
			return err
		}

		key, err := ri.mapToStructFromSlice(rawKey, keys)
		if err != nil {
			return err
		}

		keysElem.Set(reflect.Append(keysElem, key))
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func (ri DBIndex) Find(key interface{}, value interface{}) error {
	rawKey, err := ri.structToMap(key)
	if err != nil {
		return err
	}

	rawKeyBytes, err := json.Marshal(rawKey)
	if err != nil {
		return err
	}

	rows, err := ri.adapter.Find(rawKeyBytes)
	if err != nil {
		return err
	}

	if rows == nil {
		return bosherr.Errorf("Expected to find entries for key '%s'", rawKeyBytes)
	}

	if !rows.Next() {
		return bpindex.ErrNotFound
	}

	var rawEntryBytes []byte

	err = rows.Scan(&rawEntryBytes)
	if err != nil {
		return err
	}

	if rows.Next() {
		// todo should rows.Err() be called before returning
		return bosherr.Errorf("Expected to not find more than 1 entry for key '%s'", rawKeyBytes)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawEntryBytes, &value)
	if err != nil {
		return err
	}

	return nil
}

func (ri DBIndex) Save(key interface{}, value interface{}) error {
	rawKey, err := ri.structToMap(key)
	if err != nil {
		return err
	}

	rawKeyBytes, err := json.Marshal(rawKey)
	if err != nil {
		return err
	}

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	res, err := ri.adapter.Save(rawKeyBytes, valueBytes)
	if err != nil {
		return err
	}

	numRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numRows != 1 {
		return bosherr.Errorf("Expected to update 1 entry for key '%s'; updated '%d'", rawKeyBytes, numRows)
	}

	return nil
}

func (ri DBIndex) Remove(key interface{}) error {
	rawKey, err := ri.structToMap(key)
	if err != nil {
		return err
	}

	rawKeyBytes, err := json.Marshal(rawKey)
	if err != nil {
		return err
	}

	err = ri.adapter.Remove(rawKeyBytes)
	if err != nil {
		return err
	}

	return nil
}

func (ri DBIndex) FindLocked(key interface{}, value interface{}) (bpindex.LockedRecord, error) {
	rawKey, err := ri.structToMap(key)
	if err != nil {
		return dbIndexLockedRecord{}, err
	}

	rawKeyBytes, err := json.Marshal(rawKey)
	if err != nil {
		return dbIndexLockedRecord{}, err
	}

	tx, rows, err := ri.adapter.FindLocked(rawKeyBytes)
	if err != nil {
		return dbIndexLockedRecord{}, err
	}

	rec := dbIndexLockedRecord{
		rawKeyBytes: rawKeyBytes,
		notFound:    false,
		tx:          tx,
		dbIndex:     ri,

		logTag: "dbIndexLockedRecord",
		logger: ri.logger,
	}

	if rows == nil {
		// todo check error
		rec.Release()
		return dbIndexLockedRecord{}, bosherr.Errorf("Expected to find entries for key '%s'", rawKeyBytes)
	}

	if !rows.Next() {
		rec.notFound = true

		// dbIndexLockedRecord should keep lock since not finding any rows is ok
		return rec, bpindex.ErrNotFound
	}

	var rawEntryBytes []byte

	err = rows.Scan(&rawEntryBytes)
	if err != nil {
		// todo check error
		rec.Release()
		return dbIndexLockedRecord{}, err
	}

	if rows.Next() {
		// todo check error
		rec.Release()

		// todo should rows.Err() be called before returning
		return dbIndexLockedRecord{}, bosherr.Errorf("Expected to not find more than 1 entry for key '%s'", rawKeyBytes)
	}

	err = rows.Err()
	if err != nil {
		// todo check error
		rec.Release()
		return dbIndexLockedRecord{}, err
	}

	err = json.Unmarshal(rawEntryBytes, &value)
	if err != nil {
		// todo check error
		rec.Release()
		return dbIndexLockedRecord{}, err
	}

	return rec, err
}

// Release releases possibly allocated transaction
// todo refactor to not allow !nil this case
func (r dbIndexLockedRecord) Release() error {
	if r.tx != nil {
		r.dbIndex.adapter.ReleaseLocked(r.tx)
	}

	return nil
}

func (r dbIndexLockedRecord) Save(value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var res sql.Result

	if r.notFound {
		r.logger.Debug(r.logTag, "Inserting locked record")

		res, err = r.dbIndex.adapter.InsertLocked(r.tx, r.rawKeyBytes, valueBytes)
		if err != nil {
			if err == ErrExists {
				return bpindex.ErrChanged
			}

			return err
		}
	} else {
		r.logger.Debug(r.logTag, "Updating locked record")

		res, err = r.dbIndex.adapter.UpdateLocked(r.tx, r.rawKeyBytes, valueBytes)
		if err != nil {
			return err
		}
	}

	numRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numRows != 1 {
		return bosherr.Errorf("Expected to update 1 entry for key '%s'; updated '%d'", r.rawKeyBytes, numRows)
	}

	return nil
}

func (r dbIndexLockedRecord) Remove() error {
	return r.dbIndex.adapter.RemoveLocked(r.tx, r.rawKeyBytes)
}

// structToMap extracts fields from a struct and populates a map
func (ri DBIndex) structToMap(s interface{}) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	st := reflect.TypeOf(s)
	stv := reflect.ValueOf(s)

	if stv.Kind() != reflect.Struct {
		return res, bosherr.Errorf(
			"Must be reflect.Struct: %#v (%#v)", stv, ri.kindToStr(stv.Kind()))
	}

	for i := 0; i < st.NumField(); i++ {
		// Do not export private fields; private fields have PkgPath set.
		// http://golang.org/pkg/reflect/#StructField
		if st.Field(i).PkgPath == "" {
			res[st.Field(i).Name] = stv.Field(i).Interface()
		}
	}

	return res, nil
}

// mapToStruct returns new struct value with data from a map
func (ri DBIndex) mapToStruct(m map[string]interface{}, t interface{}) (reflect.Value, error) {
	return ri.mapToNewStruct(m, reflect.ValueOf(t).Elem().Type())
}

// mapToStructFromSlice returns new struct value with data from a map
func (ri DBIndex) mapToStructFromSlice(m map[string]interface{}, t interface{}) (reflect.Value, error) {
	slice := reflect.ValueOf(t).Elem()

	if slice.Kind() != reflect.Slice {
		return reflect.Value{}, bosherr.Errorf(
			"Must be reflect.Slice: %#v (%#v)",
			slice, ri.kindToStr(slice.Kind()),
		)
	}

	return ri.mapToNewStruct(m, slice.Type().Elem())
}

// mapToNewStruct returns new struct of type t with data from a map
func (ri DBIndex) mapToNewStruct(m map[string]interface{}, t reflect.Type) (reflect.Value, error) {
	if t.Kind() != reflect.Struct {
		return reflect.Value{}, bosherr.Errorf(
			"Must be reflect.Struct: %#v (%#v)",
			t, ri.kindToStr(t.Kind()),
		)
	}

	newStruct := reflect.New(t).Elem()

	for k, v := range m {
		f := newStruct.FieldByName(k)
		if f.IsValid() && f.CanSet() {
			// todo float64 -> int
			// todo pointer values
			// todo slices
			f.Set(reflect.ValueOf(v))
		}
	}

	return newStruct, nil
}

// todo consolidate with file_index
var kindToStrMap = map[reflect.Kind]string{
	reflect.Invalid:       "Invalid",
	reflect.Bool:          "Bool",
	reflect.Int:           "Int",
	reflect.Int8:          "Int8",
	reflect.Int16:         "Int16",
	reflect.Int32:         "Int32",
	reflect.Int64:         "Int64",
	reflect.Uint:          "Uint",
	reflect.Uint8:         "Uint8",
	reflect.Uint16:        "Uint16",
	reflect.Uint32:        "Uint32",
	reflect.Uint64:        "Uint64",
	reflect.Uintptr:       "Uintptr",
	reflect.Float32:       "Float32",
	reflect.Float64:       "Float64",
	reflect.Complex64:     "Complex64",
	reflect.Complex128:    "Complex128",
	reflect.Array:         "Array",
	reflect.Chan:          "Chan",
	reflect.Func:          "Func",
	reflect.Interface:     "Interface",
	reflect.Map:           "Map",
	reflect.Ptr:           "Ptr",
	reflect.Slice:         "Slice",
	reflect.String:        "String",
	reflect.Struct:        "Struct",
	reflect.UnsafePointer: "UnsafePointer",
}

func (ri DBIndex) kindToStr(k reflect.Kind) string {
	return kindToStrMap[k]
}
