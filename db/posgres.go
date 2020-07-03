package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

type Postgres struct {
	SqlDB *sql.DB
}

func NewPostgres() *Postgres {
	return &Postgres{
		SqlDB: createPGClient(),
	}
}

func createPGClient() *sql.DB {

	// Open the connection
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	//POSTGRES_URL="host=42.159.86.191 port=5432 user=46d1a69b-6cd1-4b94-b009-537e2d575bba password=ssc8u7occfhqm3q6gkhm0gvcua dbname=ecd73592-abcd-4a8e-a7c9-26e1d5bab72c parsetime=true sslmode=disable"

	if err != nil {
		glog.Error("create pg connection err:", err)
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		glog.Error("ping postgres err:", err)
		panic(err)
	}

	fmt.Println("Successfully connected postgres!")

	return db
}
