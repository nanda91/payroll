package database

import (
	"fmt"
	"math/rand"
	"payroll/domain/model"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	db.Create(&model.User{
		Username: "admin",
		Password: string(password),
		Salary:   0,
		Role:     model.RoleAdmin,
	})

	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		pass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		db.Create(&model.User{
			Username: fmt.Sprintf("employee%d", i),
			Password: string(pass),
			Salary:   float64(3_000_000 + rand.Intn(2_000_000)),
			Role:     model.RoleEmployee,
		})
	}
}
