package dao

import (
	"db/repository/db/model"
)

// migration 迁移
func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
		)

	if err != nil {
		panic(err)
	}
}
