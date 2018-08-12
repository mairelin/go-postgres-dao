package go_postgres_dao



import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"reflect"
	"strconv"
	"time"
	"strings"
	"errors"
)

type PostgresDB struct {
	*sql.DB
	ConnString string
	Driver string
}

// This method initialize a database connection, with the values passed as parameter to the parent PostgresDB
// returns *sql.DB, error
func (pdb *PostgresDB) InitDB() (*sql.DB, error) {
	var err error

	database, err := sql.Open(pdb.Driver, pdb.ConnString)

	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err == nil {
		fmt.Println("Successfully connected!")
	}

	pdb.DB = database
	return pdb.DB, err
}

// Description: Search a list of a given struct type
// Params:
//  - model: type that has been mapped before with a database table
//  - limit: This is for pagination and indicates how many rows will return per consult,
//         if value is 0, returns all rows data.
//  - offset: This is for pagination and indicates how many rows have to be avoided
//  - filters: a map that contains key values for search string, key will represent the column
//            to be filtered and the value represents the value to compare, is this comes
//            empty returns rows with out condition.
// Returns:
//  - *sql.Rows: Result list
//  - error: Error type
func (pdb PostgresDB) ListAllPaginated(model interface{}, limit int, offset int, filters map[string]interface{}) (*sql.Rows, error) {
	t := reflect.ValueOf(model).Elem()
	var resRows  *sql.Rows
	var resError error
	sqlString := " SELECT "
	sqlString = concatColumns(t, sqlString) + " FROM " + t.Type().Name() + " where deleted_at is null"

	var i int
	for k, _ := range filters {
		i++
		sqlString = sqlString + " AND " + k + " = $" + strconv.Itoa(i)
	}

	newValues := pdb.extractFilterValues(filters)

	if !(limit == 0 && offset == 0) {
		if len(newValues) == 0 {
			sqlString = sqlString + " order by id desc LIMIT $1 OFFSET $2 "
		} else {
			newValues = append(newValues, strconv.Itoa(limit))
			newValues = append(newValues, strconv.Itoa(offset))
			sqlString = sqlString + " order by id desc LIMIT $" + strconv.Itoa(len(filters)+1) + " OFFSET $" + strconv.Itoa(len(filters)+2)
		}
	}

	resRows, resError = pdb.Query(sqlString, newValues...)

	return resRows, resError
}

func (pdb PostgresDB) extractFilterValues(filters map[string]interface{}) []interface{} {
	var newValues []interface{}
	for _, v := range filters {
		newValues = append(newValues, getValue(reflect.ValueOf(v)))
	}
	return newValues
}

// Description: search a given type struct passed as parameter by id
// Params:
//   - model: type that has been mapped before with a database table
//   - id: value for the filter.
// Returns:
//  - *sql.Row: Result row
func (pdb PostgresDB) GetById(model interface{}, id uint)  *sql.Row  {
	t := reflect.ValueOf(model).Elem()
	sqlString := " SELECT "
	sqlString = concatColumns(t, sqlString) + " FROM " + t.Type().Name() + " WHERE ID = $1 and deleted_at is null;"
	return pdb.QueryRow(sqlString, id)
}

// Description: delete a given type struct passed as parameter
// Params:
//  - model: type that has been mapped before with a database table
//  - id: identification of the row to be deleted.
// Returns:
//  - sql.Result returns id of affected row
//  - error
func (pdb PostgresDB) Delete(model interface{}, id uint)  (sql.Result, error) {
	t := reflect.ValueOf(model).Elem()
	sqlString := " UPDATE " + t.Type().Name() + " SET deleted_at = to_timestamp('" + time.Now().Format("2006-01-02 15:04:05") + "', 'YYYY-MM-DD HH24:MI:SS') WHERE ID = $1;"
	return pdb.Exec(sqlString, id)
}

// Description: update a given type struct passed as parameter
// Params:
//  - model: type that has been mapped before with a database table
//  - id: identification of the row to be updated.
// Returns:
//  - sql.Result info about updated data
//  - error
func (pdb PostgresDB) Update(model interface{}, id uint)  (sql.Result, error) {
	t := reflect.ValueOf(model).Elem()
	sqlString := " UPDATE " + t.Type().Name() + " SET "
	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Name != "ID" {
			sqlString = sqlString + t.Type().Field(i).Tag.Get("model") + " = $" + strconv.Itoa(i)
			if i != (t.NumField() - 1) {
				sqlString = sqlString + ","
			}
		}
	}

	newValues := pdb.getValueParams(t)
	newValues = append(newValues, strconv.Itoa(int(id)))
	idParam := strconv.Itoa(len(newValues))

	sqlString = sqlString + ", updated_at = to_timestamp('" + time.Now().Format("2006-01-02 15:04:05") + "', 'YYYY-MM-DD HH24:MI:SS') WHERE ID = $" + idParam

	return pdb.Exec(sqlString,  newValues...)
}

// Description: insert a new type
// Params:
//  - model: type that has been mapped before with a database table
// Returns:
//  - sql.Result info about inserted data
//  - error
func (pdb PostgresDB) Create(model interface{})   error {
	t := reflect.ValueOf(model).Elem()
	sqlString := "INSERT INTO "+ t.Type().Name() + " ("

	//Get fields
	sqlString = concatColumns(t, sqlString)

	//Get values
	sqlString = sqlString + " ) VALUES ( nextval('" + t.Type().Name() + "_seq'), "
	sqlString = concatValues(t, sqlString)

	newValues := pdb.getValueParams(t)

	_, err := pdb.Exec(sqlString, newValues...)
	fmt.Errorf("Error executing statement", err)

	return err
}

func (pdb PostgresDB) getValueParams(t reflect.Value) []interface{} {
	var newValues []interface{}
	for i := 0; i < t.NumField(); i++ {
		if t.Type().Field(i).Name != "ID" {
			newValues = append(newValues, getValue(t.Field(i)))
		}
	}
	return newValues
}

// Description: create value sqlString from fields of a given  value
// Params:
//  - valueOf: Value for a reflected type
//  - sqlString: String to be concatenated
// Returns:
//  - string: string with the format (val, val2, val3...)
func concatValues(valueOf reflect.Value, sqlString string) string {
	for i := 0; i < valueOf.NumField(); i++ {

		if valueOf.Type().Field(i).Name != "ID" {
			sqlString = sqlString +  "$" + strconv.Itoa(i)

			if i != (valueOf.NumField() - 1) {
				sqlString = sqlString + ","
			} else {
				sqlString = sqlString + ");"
			}
		}
	}
	return sqlString
}


// Description: create columns sqlString from fields of a given value tags,
//              this method search for the tag 'model' added to fields
//              of the current struct
// Params:
//  - valueOf: Value for a reflected type
//  - sqlString: String to be concatenated
// Returns:
//  - string: string with the format (column1, column2, column...), those values
func concatColumns(valueOf reflect.Value, sqlString string) string {
	for i := 0; i < valueOf.NumField(); i++ {
		if len(valueOf.Type().Field(i).Tag.Get("model")) > 0 {
			sqlString = sqlString + valueOf.Type().Field(i).Tag.Get("model")
			if i != (valueOf.NumField() - 1) {
				sqlString = sqlString + ","
			}
		}
	}
	return sqlString
}


// Description: return the value string as sql type of a given value
// Params:
//  - fieldValue: field value to be evaluated
// Returns:
//  - string: string with the value ready to concat to sql statement
func getValue(fieldValue reflect.Value) string {
	res := ""
	val := fieldValue.Type().Name()
	switch val {
	case "int":
		res = strconv.FormatInt(fieldValue.Int(), 10)
	case "uint":
		res = strconv.FormatUint(fieldValue.Uint(), 10)
	case "string":
		res =  fieldValue.String()
	case "bool":
		res =  strconv.FormatBool(fieldValue.Bool())
	case "Time":
		if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
			res =  "to_timestamp('" + fieldValue.Interface().(time.Time).Format("2006-01-02 15:04:05") + "', 'YYYY-MM-DD HH24:MI:SS')"		// ErrorR in .Interface()
		}
	}
	return res
}

// Description: create a table with the tags mapped with in the struct passed as parameter
// Params:
//  - model: interface for table creating, every field that will to be mapped
//           must have the tags:
//                              'model' with the name of the column
//                              'type' with type for the column on the database
//           can also have:
//                             'mandatory' for not null values
//                             'unique' for add unique constraint
//                             'reference' for add foreign
//           Example of a mapped struct:
//								type NoteRating struct {
//									ID uint   `model:"id" type:"bigserial" constraint:"rating_Note_pk PRIMARY KEY(id)" `
//									NoteID uint `model:"Note_id" type:"bigserial" mandatory:"true" reference:"Note(ID)"`
//									Rating int `model:"rating" type:"smallint" mandatory:"true"`
//									Date time.Time `model:"rating_date" type:"timestamp" mandatory:"true"`
//								}
// Returns:
//  - error
// Note: is mandatory add the field 'ID uint'
func (pdb PostgresDB) CreateTable(model interface{}) error {
	structFieldValue := reflect.ValueOf(model).Elem()

	query := "CREATE TABLE " + structFieldValue.Type().Name() + " ( "

	for i := 0; i < structFieldValue.NumField() ; i++ {

		if len(structFieldValue.Type().Field(i).Tag.Get("model")) > 0 && len(structFieldValue.Type().Field(i).Tag.Get("type"))> 0{

			query = query + structFieldValue.Type().Field(i).Tag.Get("model") + " " + structFieldValue.Type().Field(i).Tag.Get("type")

			if len(structFieldValue.Type().Field(i).Tag.Get("mandatory")) > 0 {
				query = query +  " NOT NULL "
			}

			if len(structFieldValue.Type().Field(i).Tag.Get("unique")) > 0 {
				query = query +  " UNIQUE "
			}


			if len(structFieldValue.Type().Field(i).Tag.Get("constraint")) > 0 {
				query = query + ", CONSTRAINT " +  structFieldValue.Type().Field(i).Tag.Get("constraint")
			} else  if len(structFieldValue.Type().Field(i).Tag.Get("reference")) > 0 {
				query = query + " references " +  structFieldValue.Type().Field(i).Tag.Get("reference")
			}

			if i != (structFieldValue.NumField() - 1) {
				query = query + ","
			} else {
				query = query + ", deleted_at timestamp, updated_at timestamp);"
			}
		}
	}
	_, err := pdb.Exec(query)
	fmt.Errorf("Error executing query", err)
	return err
}

// Description: create a sequence for a given struct with the name {structName}_seq
// Params:
//  - model: struct for the sequence creation
// Returns:
//  - error
func (pdb PostgresDB) CreateSequence(model interface{}) error {
	t := reflect.TypeOf(model)
	query := "CREATE SEQUENCE " + t.Name() + "_seq" + " increment 1 minvalue 1 maxvalue 9223372036854775807  start 1 cache 1;"
	_, err := pdb.Exec(query)
	fmt.Errorf("Error executing statement", err)
	return err
}

// Description: drop sequence for a given struct with the name {structName}_seq
// Params:
//  - model: struct for the sequence creation
// Returns:
//  - error
func (pdb PostgresDB) DropSequence(model interface{}) error {
	t := reflect.TypeOf(model)
	query := "DROP SEQUENCE " + t.Name() + "_seq ;"
	_, err := pdb.Exec(query)
	fmt.Errorf("Error executing statement", err)
	return err
}

// Description: drop table for a given struct with the name
// Params:
//  - model: struct for the sequence creation
// Returns:
//  - error
func (pdb PostgresDB) DropTable(model interface{}) error {
	t := reflect.TypeOf(model)
	query := "DROP TABLE " + t.Name()
	_, err := pdb.Exec(query)
	fmt.Errorf("Error executing statement", err)
	return err
}

// Description: check if a table for a given struct exists at the database
// Params:
//  - model: struct for the sequence creation
// Returns:
//  - error
func (pdb PostgresDB) CheckIfExists(model interface{}) (bool, error) {
	var res error
	exists := true
	if pdb.DB != nil {
		t := reflect.TypeOf(model)
		sqlString := " SELECT * FROM " + t.Name() + " limit 1;"
		_, err := pdb.Query(sqlString)
		if err != nil && strings.Contains(err.Error(), "does not exist") {
			exists = false
		}
	} else {
		fmt.Errorf("Database is not initialized")
		res = errors.New("Database is not initialized")
	}
	return exists, res
}
