package orm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	// TimeFormat for datatime
	TimeFormat = "2006-01-02 15:04:05.999999999"
	// ErrNotValid means database has NOT beed configured.
	ErrNotValid = errors.New("NOT VALID")
	// ErrNotFound means record can NOT be found in database.
	ErrNotFound = errors.New("RECORD NOT FOUND")
	// ErrUnsupportType means using type that unsupported
	ErrUnsupportType = errors.New("UNSUPPORT TYPE")
	// ErrBadParam means parameter is in bad format.
	ErrBadParam = errors.New("BAD PARAMETER")
)

var db *sql.DB

// CreateTable by value type.
func CreateTable(v interface{}) error {
	if db == nil {
		return ErrNotValid
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrUnsupportType
	}

	de := rv.Elem()
	dt := de.Type()
	if de.Kind() != reflect.Struct {
		return ErrUnsupportType
	}

	var builder strings.Builder
	builder.WriteString("CREATE TABLE IF NOT EXISTS `")
	builder.WriteString(strings.ToLower(dt.Name()))
	builder.WriteString("`(\n")

	fields := []string{}
	hasID := false

	for i := 0; i < dt.NumField(); i++ {
		fv := de.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		ft := dt.Field(i)
		tag := ft.Tag.Get("mysql")
		if tag == "-" || tag == "" {
			continue
		}

		name := strings.ToLower(ft.Name)
		if name == "id" {
			fields = append(fields, "`id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY")
			hasID = true
		} else {
			fields = append(fields, fmt.Sprintf("`%s` %s", name, tag))
		}
	}

	if len(fields) == 0 {
		return ErrBadParam
	}

	if !hasID {
		builder.WriteString("`id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,\n")
	}

	builder.WriteString(strings.Join(fields, ",\n"))
	builder.WriteString(") DEFAULT CHARSET utf8;")

	_, err := db.Exec(builder.String())
	return err
}

// ConnectDB starts connection with MySQL server
func ConnectDB(addr, user, pswd, dbName string, maxConns int) error {
	param := mysql.NewConfig()
	param.Net = "tcp"
	param.Addr = addr
	param.User = user
	param.Passwd = pswd
	param.DBName = dbName
	param.MultiStatements = true
	param.Params = map[string]string{
		"charset":   "utf8",
		"collation": "utf8_general_ci",
	}

	conn, err := sql.Open("mysql", param.FormatDSN())
	if err != nil {
		return err
	}

	conn.SetMaxOpenConns(maxConns)

	db = conn
	return nil
}

// Exec SQL without results like DELETE/UPDATE/INSERT
func Exec(sql string, args ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, fmt.Errorf("orm.Exec on invalid connection")
	}

	return db.Exec(sql, args...)
}

// Query with results.
func Query(sql string, args ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, fmt.Errorf("orm.Exec on invalid connection")
	}

	return db.Query(sql, args...)
}

// Scan current row in result set into struct.
func Scan(rows *sql.Rows, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrUnsupportType
	}

	de := rv.Elem()
	if de.Kind() != reflect.Struct {
		return rows.Scan(v)
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	colSize := len(cols)
	vals := make([]interface{}, colSize, colSize)
	ptrs := make([][]byte, colSize, colSize)

	for i := 0; i < colSize; i++ {
		ptrs[i] = make([]byte, 1, 1)
		vals[i] = &ptrs[i]
	}

	rows.Scan(vals...)

	for i := 0; i < colSize; i++ {
		fv := de.FieldByNameFunc(func(name string) bool { return strings.ToLower(name) == cols[i] })
		if err = deserialize(fv, ptrs[i]); err != nil {
			return err
		}
	}

	return nil
}

// Insert data into database.
func Insert(v interface{}) (sql.Result, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, ErrUnsupportType
	}

	de := rv.Elem()
	dt := de.Type()
	if de.Kind() != reflect.Struct {
		return nil, ErrUnsupportType
	}

	var builder strings.Builder

	builder.WriteString("INSERT INTO `")
	builder.WriteString(strings.ToLower(dt.Name()))
	builder.WriteString("`(")

	keys := []string{}
	holders := []string{}
	vals := []interface{}{}

	for i := 0; i < dt.NumField(); i++ {
		ft := dt.Field(i)
		name := strings.ToLower(ft.Name)
		tag := ft.Tag.Get("mysql")
		if tag == "-" || tag == "" || name == "id" {
			continue
		}

		fv := de.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		val, err := serialize(fv)
		if err != nil {
			return nil, err
		}

		keys = append(keys, "`"+name+"`")
		holders = append(holders, "?")
		vals = append(vals, val)
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("orm.Insert No valid fields found in record: %+v", v)
	}

	builder.WriteString(strings.Join(keys, ","))
	builder.WriteString(") VALUES(")
	builder.WriteString(strings.Join(holders, ","))
	builder.WriteString(");")

	return Exec(builder.String(), vals...)
}

// Read one record from database
func Read(v interface{}, cols ...string) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrUnsupportType
	}

	de := rv.Elem()
	dt := de.Type()
	if de.Kind() != reflect.Struct {
		return ErrUnsupportType
	}

	var builder strings.Builder

	builder.WriteString("SELECT * FROM `")
	builder.WriteString(strings.ToLower(dt.Name()))
	builder.WriteString("` WHERE ")

	conditions := []string{}
	vals := []interface{}{}

	if cols != nil && len(cols) > 0 {
		for _, col := range cols {
			fv := de.FieldByNameFunc(func(name string) bool { return strings.ToLower(name) == col })
			if !fv.IsValid() {
				return ErrBadParam
			}

			val, err := serialize(fv)
			if err != nil {
				return err
			}

			conditions = append(conditions, "`"+col+"`=?")
			vals = append(vals, val)
		}
	} else {
		fv := de.FieldByNameFunc(func(name string) bool { return strings.ToLower(name) == "id" })
		if !fv.IsValid() || fv.Kind() != reflect.Int64 {
			return ErrBadParam
		}

		conditions = append(conditions, "`id`=?")
		vals = append(vals, fv.Int())
	}

	builder.WriteString(strings.Join(conditions, " AND "))

	rows, err := Query(builder.String(), vals...)
	if err != nil {
		return err
	}

	defer rows.Close()

	if !rows.Next() {
		return ErrNotFound
	}

	return Scan(rows, v)
}

// Update one record from database
func Update(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrUnsupportType
	}

	de := rv.Elem()
	dt := de.Type()
	if de.Kind() != reflect.Struct {
		return ErrUnsupportType
	}

	var builder strings.Builder

	builder.WriteString("UPDATE `")
	builder.WriteString(strings.ToLower(dt.Name()))
	builder.WriteString("` SET ")

	id := int64(-1)
	keys := []string{}
	vals := []interface{}{}

	for i := 0; i < dt.NumField(); i++ {
		ft := dt.Field(i)
		name := strings.ToLower(ft.Name)
		tag := ft.Tag.Get("mysql")
		if tag == "-" || tag == "" {
			continue
		}

		fv := de.Field(i)
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		if name == "id" {
			id = fv.Int()
		} else {
			val, err := serialize(fv)
			if err != nil {
				return err
			}

			keys = append(keys, "`"+name+"`=?")
			vals = append(vals, val)
		}
	}

	if id < 0 {
		return ErrBadParam
	}

	builder.WriteString(strings.Join(keys, ","))
	builder.WriteString(fmt.Sprintf(" WHERE `id`=%d", id))

	_, err := Exec(builder.String(), vals...)
	if err != nil {
		return err
	}

	return nil
}

// Delete a record from data by ID
func Delete(table string, id int64) error {
	_, err := Exec("DELETE FROM `"+table+"` WHERE `id`=?", id)
	return err
}

func serialize(v reflect.Value) (interface{}, error) {
	switch v.Kind() {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Array, reflect.Slice:
		data, err := json.Marshal(v.Interface())
		if err != nil {
			return nil, err
		}

		return string(data), nil
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t := v.Interface().(time.Time)
			return t.Format(TimeFormat), nil
		}

		data, err := json.Marshal(v.Interface())
		if err != nil {
			return nil, err
		}

		return string(data), nil
	default:
		return nil, ErrUnsupportType
	}
}

func deserialize(v reflect.Value, raw []byte) error {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		n, err := strconv.Atoi(string(raw))
		if err != nil {
			return err
		}

		v.SetBool(n != 0)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(string(raw), 10, 64)
		if err != nil {
			return err
		}

		v.SetInt(n)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(string(raw), 10, 64)
		if err != nil {
			return err
		}

		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(string(raw), 64)
		if err != nil {
			return err
		}

		v.SetFloat(n)
	case reflect.String:
		v.SetString(string(raw))
	case reflect.Array, reflect.Slice:
		err := json.Unmarshal(raw, v.Addr().Interface())
		if err != nil {
			return err
		}
	case reflect.Struct:
		if v.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(TimeFormat, string(raw))
			if err != nil {
				return err
			}

			v.Set(reflect.ValueOf(t))
			return nil
		}

		err := json.Unmarshal(raw, v.Addr().Interface())
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportType
	}

	return nil
}
