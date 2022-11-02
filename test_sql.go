package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const dsn = "dump:111@tcp(192.168.157.128:6001)/mo?charset=utf8mb4"
const driverName = "mysql"

func test1() {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	if err != nil {
		panic(err.Error())
	}

	// Prepare statement for drop table
	stmtDrop, err := db.Prepare("drop database if exists gosqldriver")
	if err != nil {
		panic(err.Error())
	}
	defer stmtDrop.Close() // Close the statement when we leave main() / the program terminates
	stmtDrop.Exec()

	// Prepare statement for create database
	stmtCreateDatabase, err := db.Prepare("create database gosqldriver") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	defer stmtCreateDatabase.Close() // Close the statement when we leave main() / the program terminates
	stmtCreateDatabase.Exec()

	_, err = db.Exec("use gosqldriver") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}

	// Prepare statement for create table
	stmtCreate, err := db.Prepare("create table t1(c1 int, c2 int)") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	defer stmtCreate.Close()
	stmtCreate.Exec()

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO t1 VALUES( ?, ? )") // ? = placeholder
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT c2 FROM t1 WHERE c1 = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error())
		}
	}

	squareNum := 0

	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error())
	}
	if squareNum != 169 {
		fmt.Println("failed!")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err.Error())
	}
	stmtdelete, err := tx.Prepare("delete from t1 where c1>0;")
	if err != nil {
		panic(err.Error())
	}
	stmtdelete.Exec()
	rows, _ := tx.Query("select count(*) from t1")
	var count int
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count != 1 {
		fmt.Println("test failed!")
		return
	}

	err = tx.Rollback()
	if err != nil {
		panic(err.Error())
	}

	tx, err = db.Begin()
	if err != nil {
		panic(err.Error())
	}

	err = tx.QueryRow("select count(*) from t1").Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count != 25 {
		fmt.Println("test failed!")
		return
	}

	stmtdelete, err = tx.Prepare("delete from t1 where c1>0;")
	if err != nil {
		panic(err.Error())
	}
	stmtdelete.Exec()
	tx.Commit()

	stmtResult, err := db.Prepare("select count(*) from t1")
	if err != nil {
		panic(err.Error())
	}
	err = stmtResult.QueryRow().Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	if count != 1 {
		fmt.Println("test failed!")
		return
	}

}

func test2() {

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxIdleTime(100)
	db.SetConnMaxLifetime(100)

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	stats := db.Stats()
	fmt.Printf("dbstats: \nMaxOpenConnections %v\nIdle %v\nOpenConnections %v\nInUse %v\nWaitCount %v\nWaitDuration %v\nMaxIdleClosed %v\nMaxIdleTimeClosed %v\nMaxLifetimeClosed %v\n\n", stats.MaxOpenConnections, stats.Idle, stats.OpenConnections, stats.InUse, stats.WaitCount, stats.WaitDuration, stats.MaxIdleClosed, stats.MaxIdleTimeClosed, stats.MaxLifetimeClosed)
}

func main() {
	test1()
	test2()
	fmt.Println("test success!")
}
