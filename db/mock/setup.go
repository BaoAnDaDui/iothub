package mock

import (
	"log"

	"iothub/db/mysql"
	"iothub/db/sqlite"

	"gorm.io/gorm"
)

func NewSqliteConnTest() *gorm.DB {
	cfg := sqlite.Config{
		FilePath: ":memory:",
		//FilePath: "./test.repo",
	}
	conn, err := sqlite.Connect(cfg)
	if err != nil {
		log.Fatalf("get db conn error: %v", err)
	}
	return conn
}

func NewMySqlConnTest() *gorm.DB {
	cfg := mysql.Config{
		Host:            "localhost",
		Port:            "3306",
		User:            "root",
		Password:        "123",
		DB:              "iothub",
		Charset:         "utf8",
		Timezone:        "Asia%2FShanghai",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30,
	}
	conn, err := mysql.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
