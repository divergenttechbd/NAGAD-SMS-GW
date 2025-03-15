package seeders

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"myproject/models"
	"myproject/utils"
)

func SeedData() {
	db := utils.GetDB()

	// Check if the seeder has already been run
	var count int64
	if err := db.Model(&models.SeederLog{}).Where("seeder_name = ?", "initial_seed").Count(&count).Error; err != nil {
		log.Fatalf("Failed to check seeder logs: %v", err)
	}

	if count > 0 {
		log.Println("Seeder has already been run. Skipping...")
		return
	}

	// Add a role
	adminRole := models.Role{
		Name: "admin",
	}
	if err := db.Create(&adminRole).Error; err != nil {
		log.Fatalf("Failed to seed role: %v", err)
	}

	// Add a user
	hashedPassword := md5.Sum([]byte("12345")) // Hash the password using MD5
	adminUser := models.User{
		Username: "admin",
		Email:    "admin@nagad.com",
		Password: hex.EncodeToString(hashedPassword[:]), // Convert the hash to a hex string
		Roles:    []models.Role{adminRole},              // Assign the admin role to the user
	}

	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatalf("Failed to seed user: %v", err)
	}

	// Add permissions
	permissions := []models.Permission{
		{Name: "create_user"},
		{Name: "edit_user"},
		{Name: "delete_user"},
		{Name: "view_user"},
	}

	for _, permission := range permissions {
		if err := db.Create(&permission).Error; err != nil {
			log.Fatalf("Failed to seed permission: %v", err)
		}
	}

	// Assign permissions to the admin role
	if err := db.Model(&adminRole).Association("Permissions").Append(permissions); err != nil {
		log.Fatalf("Failed to assign permissions to admin role: %v", err)
	}

	// Log the seeder execution
	seederLog := models.SeederLog{
		SeederName: "initial_seed",
	}
	if err := db.Create(&seederLog).Error; err != nil {
		log.Fatalf("Failed to log seeder execution: %v", err)
	}

	log.Println("Seed data added successfully!")
}
