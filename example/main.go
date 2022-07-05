package main

import (
	"log"

	filter "github.com/pandalzy/gorm-filter"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Name     string
	Age      int
	Email    string
}

type UserFilter struct {
	Username *string       `form:"username" filter:"field:username;expr:contains"`
	Name     string        `form:"name" filter:"field:name;expr:contains"`
	Age      []interface{} `form:"age" filter:"field:age;expr:in"`
	Email    string        `form:"email" filter:"field:email"`
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm-filter?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//if err := db.AutoMigrate(&User{}); err != nil {
	//	log.Fatalln("migrate error")
	//}
	username := "a"
	uf := UserFilter{
		Username: &username,
		Name:     "l",
		Age:      []interface{}{20, 19},
		Email:    "",
	}
	var users []User
	if db, err := filter.Query(db.Model(&User{}).Debug(), &uf); err != nil {
		log.Println(err)
	} else {
		db.Find(&users)
	}
	log.Println(users)
}
