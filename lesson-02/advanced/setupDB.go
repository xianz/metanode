package main

import (
	"lesson02/utils"
	"log"
)

func main() {
	aSetupDB()
}

func aSetupDB() {
	db := utils.NewDB("setup.db")

	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}

	// Clear the users table before inserting new data
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Fatalf("clear users table: %v", err)
	}

	// Prepare test data: a slice of User structs to be inserted
	users := []User{
		{Name: "Alice", Email: "alice@example.com", Age: 28, Status: "active"},
		{Name: "Alice1", Email: "alice1@example.com", Age: 28, Status: "active"},
		{Name: "Alice2", Email: "alice2@example.com", Age: 28, Status: "inactive"},
		{Name: "Alice3", Email: "alice3@example.com", Age: 28, Status: "inactive"},
		{Name: "Bob", Email: "bob@example.com", Age: 31, Status: "active"},
		{Name: "Bob1", Email: "bob1@example.com", Age: 32, Status: "active"},
		{Name: "Bob2", Email: "bob2@example.com", Age: 33, Status: "inactive"},
		{Name: "Bob3", Email: "bob3example.com", Age: 34, Status: "inactive"},
	}

	if err := db.CreateInBatches(users, 3).Error; err != nil {
		log.Fatalf("seed users: %v", err)
	}

	var count int64
	if err := db.Model(&User{}).Count(&count).Error; err != nil {
		log.Fatalf("count users: %v", err)
	}

	if count != int64(len(users)) {
		log.Fatalf("expected %d users, got %d", len(users), count)
	}

	log.Fatalf("created %d users", count)
}
