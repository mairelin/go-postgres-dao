package go_postgres_dao



import (
	"testing"
	"fmt"
)

type TestModel struct {
	ID uint   `model:"id" type:"bigserial" constraint:"user_pk PRIMARY KEY (id)" `
	Name string `model:"name" type:"varchar(100)" mandatory:"true" unique:"true"`
}


func TestUtilAll(t *testing.T) {
	t.Run("createTable", TestCreateTable)
	t.Run("createSequence", TestCreateSeq)
	t.Run("create", TestCreate)
	t.Run("findById", TestGetById)
	t.Run("update", TestUpdate)
	t.Run("delete", TestDelete)
	t.Run("dropTable", TestDropTable)
	t.Run("dropSequence", TestDropSequence)
}

func getDataBase() PostgresDB {
	connString := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "dbname",
		"dbname", "dbname")
	db :=  PostgresDB{ConnString:connString, Driver: "postgres"}
	db.InitDB()
	return db
}

func TestCreate(t *testing.T) {
	body := TestModel{Name:"only test"}
	db  := getDataBase()
	err := db.Create(&body)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}

func TestCreateTable(t *testing.T) {
	db := getDataBase()
	err := db.CreateTable(&TestModel{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
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

func TestUpdate(t *testing.T) {
	db := getDataBase()
	res , err := db.Update(&TestModel{Name:"Nuevo Nombre"}, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}

func TestDelete(t *testing.T) {
	db := getDataBase()
	res, err := db.Delete(&TestModel{}, 1)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
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

func TestGetById(t *testing.T) {
	db := getDataBase()
	row := db.GetById(&TestModel{}, 1)
	if row == nil {
		t.Error("not found")
	} else {
		t.Log(row)
	}
}

func TestDataBaseConn(t *testing.T) {
	connString := fmt.Sprintf("host=%s port=%s user=%s  password=%s dbname=%s sslmode=disable",
		"localhost", "5432", "pecuniaapi",
		"pecuniaapi", "pecuniaapi")
	db :=  PostgresDB{ConnString:connString, Driver: "postgres"}
	_, err := db.InitDB()
	if err != nil {
		t.Error(err)
	}
	db.Close()
}
