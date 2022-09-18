package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var dataSourceName = checkHost()

func checkHost() string {
	var dbPassword, dbHost, dbUser, dbPort string
	if os.Getenv("DB_HOST") == "" {
		dbPassword = "password"
		dbHost = "127.0.0.1"
		dbUser = "root"
		dbPort = "3306"
	} else {
		dbPassword = os.Getenv("DB_PASSWORD")
		dbHost = os.Getenv("DB_HOST")
		dbUser = os.Getenv("DB_USER")
		dbPort = os.Getenv("DB_PORT")
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
}

func CreateDatabase(name string) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic()
		}
	}(db)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		panic(err)
	}
}

func CreateTable(dbName, tableName string) *sql.DB {
	db, err := sql.Open("mysql", dataSourceName+dbName)
	if err != nil {
		log.Panic(err)
	} else if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(chatID bigint primary key, length int,"+
		" changeLength bool, number bool, upperCase bool, lowerCase bool, specialCase bool)", tableName)
	_, err = db.Exec(query)
	if err != nil {
		log.Panic(err)
	}
	return db
}

func InsertData(dbName, tableName string, chatID int, pass *PasswordParam, changeLength bool) error {
	db, err := sql.Open("mysql", dataSourceName+dbName)
	if err != nil {
		log.Fatal(err)
	} else if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)
	query := fmt.Sprintf("SELECT * FROM %s WHERE chatID = %d", tableName, chatID)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	if rows.Next() {
		query = fmt.Sprintf("UPDATE %s SET length = %d, changeLength = %t, number = %t, upperCase = %t, "+
			"lowerCase = %t, specialCase = %t WHERE chatID = %d", tableName, pass.length, changeLength, pass.number, pass.upperCase,
			pass.lowerCase, pass.specialCase, chatID)
		_, err = db.Exec(query)
	} else {
		query = fmt.Sprintf("INSERT INTO %s(chatID, length, changeLength, number, upperCase,"+
			" lowerCase, specialCase) VALUES (%d,%d,%t,%t,%t,%t,%t)", tableName, chatID, pass.length, changeLength, pass.number, pass.upperCase,
			pass.lowerCase, pass.specialCase)
		_, err = db.Exec(query)
	}
	return err
}

func GetData(dbName, tableName string, chatID int) (*PasswordParam, bool, error) {
	db, err := sql.Open("mysql", dataSourceName+dbName)
	if err != nil {
		log.Fatal(err)
	} else if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	queryRow := fmt.Sprintf("select * from %s where chatID = %d", tableName, chatID)
	row := db.QueryRow(queryRow)
	var length int
	var changeLength, number, upperCase, lowerCase, specialCase bool
	err = row.Scan(&chatID, &length, &changeLength, &number, &upperCase, &lowerCase, &specialCase)
	if err != nil {
		log.Println("No data about user")
		return NewPasswordParam(20, true, true, true, false), false, fmt.Errorf("no data/got error")
	}
	return NewPasswordParam(length, number, upperCase, lowerCase, specialCase), changeLength, nil
}
