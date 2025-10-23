package main

// seedDatabase populates the database with example data using GORM
func seedDatabase() error {
	users := []User{
		{Name: "João Silva", Email: "joao@email.com", Status: "active", Age: 25},
		{Name: "Maria Santos", Email: "maria@email.com", Status: "active", Age: 30},
		{Name: "Pedro Costa", Email: "pedro@email.com", Status: "inactive", Age: 22},
		{Name: "Ana Oliveira", Email: "ana@email.com", Status: "active", Age: 28},
		{Name: "Carlos Souza", Email: "carlos@email.com", Status: "pending", Age: 35},
		{Name: "Juliana Lima", Email: "juliana@email.com", Status: "active", Age: 27},
		{Name: "Roberto Alves", Email: "roberto@email.com", Status: "inactive", Age: 19},
		{Name: "Fernanda Dias", Email: "fernanda@email.com", Status: "active", Age: 32},
	}

	// Use GORM Create to insert all users
	result := db.Create(&users)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
