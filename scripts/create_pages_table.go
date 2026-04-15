package main

import (
	"log"
	"photoset/internal/database"
	"photoset/internal/domain"
)

func main() {
	db := database.GetMySQL()
	if err := db.AutoMigrate(&domain.Page{}); err != nil {
		log.Fatal("AutoMigrate pages 失败:", err)
	}
	log.Println("pages 表创建/更新完成")
}