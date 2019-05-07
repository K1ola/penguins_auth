package database

import (
	"database/sql"
	"main/helpers"
	"main/models"
)

func RowsToUsers(rows *sql.Rows) []models.User {
	users := []models.User{}
	defer rows.Close()
	for rows.Next() {
		entry := models.User{}
		if err := rows.Scan(&entry.Login, &entry.Score, &entry.Email); err == nil {
			helpers.LogMsg(err)
		}
		users = append(users, entry)
	}
	return users
}
