/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Lite *gorm.DB
)

func SQLiteConnect() {
	var err error
	Lite, err = gorm.Open(sqlite.Open("nas.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to sqlite3 database")
	}
}
