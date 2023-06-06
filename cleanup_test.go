package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"testing"
)

//type AuthorDet struct {
//	gorm.Model
//	AuthorID string
//	Name     string
//}

func TestCleanup(t *testing.T) {
	db, _ := gorm.Open("sqlite3", "test.db")
	defer db.Close()

	cleaner := DeleteEntitesCreated(db)
	defer cleaner()

	var authors = []*Author{
		{AuthorID: "TA1", Name: "Sample Author 1"},
		{AuthorID: "TA2", Name: "Sample Author 2"},
		{AuthorID: "TA3", Name: "Sample Author 3"},
	}

	//db.Create(&authors)
	for _, author := range authors {
		db.Create(&author)
	}
}

func DeleteEntitesCreated(db *gorm.DB) func() {
	type entity struct {
		table   string
		keyName string
		key     interface{}
	}

	var entries []entity
	nameHook := "hookCleanUp"

	db.Callback().Create().After("gorm:create").Register(nameHook, func(scope *gorm.Scope) {
		fmt.Printf("Inserted Entites of %s with %s = %v\n", scope.TableName(), scope.PrimaryKey(), scope.PrimaryKeyValue())
		entries = append(entries, entity{table: scope.TableName(), keyName: scope.PrimaryKey(), key: scope.PrimaryKeyValue()})
	})

	return func() {
		defer db.Callback().Create().Remove(nameHook)
		_, transactionIn := db.CommonDB().(*sql.Tx)

		tx := db
		if !transactionIn {
			tx = db.Begin()
		}

		for i := len(entries) - 1; i >= 0; i-- {
			entry := entries[i]
			fmt.Printf("Deleting entities from %s table with key %v\n", entry.table, entry.key)
			tx.Table(entry.table).Where(entry.keyName+" = ?", entry.key).Delete("")
		}

		if !transactionIn {
			tx.Commit()
		}
	}
}
