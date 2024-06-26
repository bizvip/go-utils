/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
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
