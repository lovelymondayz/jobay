package seeder

import "gorm.io/gorm"


func SeedData(db *gorm.DB) {
	SeedSites(db)
	SeedRole(db)
	SeedClients(db)
	SeedUsers(db)
}