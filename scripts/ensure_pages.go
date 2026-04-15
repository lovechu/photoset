package main

import (
	"log"
	"photoset/internal/database"
	"photoset/internal/domain"
	"photoset/internal/repository"
)

func main() {
	db := database.GetMySQL()
	repo := repository.NewPageRepository(db)
	// Find first admin user
	var admin domain.User
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		log.Fatal("没有找到 admin 用户:", err)
	}
	err := repo.EnsureDefaultPages(admin.ID)
	if err != nil {
		log.Fatal("创建默认页面失败:", err)
	}
	log.Println("确保默认页面创建成功")
}