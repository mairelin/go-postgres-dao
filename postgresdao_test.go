package go_postgres_dao

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type TestModel struct {
	ID uint   `model:"id" type:"bigserial" constraint:"test_pk PRIMARY KEY (id)" `
	Name string `model:"name" type:"varchar(100)" mandatory:"true" `
}


func TestUtilAll(t *testing.T) {
	t.Run("createTable", TestCreateTable)
	t.Run("createSequence", TestCreateSeq)
	t.Run("dropTable", TestDropTable)
	t.Run("dropSequence", TestDropSequence)
}


func TestCRUD(t *testing.T) {
	t.Run("create", TestCreate)
	t.Run("findById", TestGetById)
	t.Run("findById", TestGetList)
	t.Run("update", TestUpdate)
}


func getDataBase() PostgresDB {
	connString := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "dbname",
		"dbuser", "dbpass")
	db :=  PostgresDB{ConnString:connString, Driver: "postgres"}
	db.InitDB()
	return db
}

func TestCreateTable(t *testing.T) {
	db := getDataBase()
	err := db.CreateTable(&TestModel{})
	if err != nil {
		t.Error(err)
	} else {
		res, _ := db.CheckIfExists(TestModel{})
		assert.True(t, res, "Table creation Failure!")
	}
}

func TestCreateSeq(t *testing.T) {
	db := getDataBase()
	err := db.CreateSequence(TestModel{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}


func TestCreate(t *testing.T) {
	body := TestModel{Name:"Mairelin Mairelin"}
	db  := getDataBase()
	err := db.Create(&body)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}

func TestUpdate(t *testing.T) {
	db := getDataBase()
	_ , err := db.Update(&TestModel{Name:"Mai"}, 1)
	if err != nil {
		t.Error(err)
	} else {
		row := db.GetById(&TestModel{}, 1)
		var res TestModel
		row.Scan(&res.ID, &res.Name)
		assert.Equal(t, "Mai",  res.Name, "Update Failure!")
	}
}

func TestGetById(t *testing.T) {
	db := getDataBase()
	row := db.GetById(&TestModel{}, 1)
	if row == nil {
		t.Error("not found")
	} else {
		var res TestModel
		row.Scan(&res.ID, &res.Name)
		assert.Equal(t, uint(1),  res.ID,"Query failure!")
	}
}

func TestGetList(t *testing.T) {
	filters := make(map[string]interface{})
	filters["name"] = "Mai"

	db := getDataBase()
	rows, err := db.ListAllPaginated(&TestModel{}, 1, 0 , filters)
	if rows == nil || err != nil {
		t.Error("not found", err.Error())
	} else {
		assert.True(t, rows.Next(), "Query failure!")
	}
}


func TestDelete(t *testing.T) {
	db := getDataBase()
	_, err := db.Delete(&TestModel{}, 1)
	if err != nil {
		t.Error(err)
	} else {
		row := db.GetById(&TestModel{}, 1)
		var res TestModel
		row.Scan( &res.Name)
		assert.Empty(t, res.Name,  "Delete failure!")
	}
}


func TestDataBaseConn(t *testing.T) {
	connString := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "dbname",
		"dbuser", "dbpass")
	db :=  PostgresDB{ConnString:connString, Driver: "postgres"}
	_, err := db.InitDB()
	if err != nil {
		t.Error(err)
	}
	db.Close()
}


func TestDropSequence(t *testing.T) {
	db := getDataBase()
	err := db.DropSequence(TestModel{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}

func TestDropTable(t *testing.T) {
	db := getDataBase()
	err := db.DropTable(TestModel{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}