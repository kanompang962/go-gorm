package main

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string
	Author      string
	Description string
	Price       uint
}

func createBook(db *gorm.DB, book *Book) error {
	result := db.Create(book)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func getBooks(db *gorm.DB) []Book {
	var books []Book
	result := db.Find(&books)
	if result.Error != nil {
		log.Fatalf("Error finding book: %v", result.Error)
	}
	return books
}

func getBook(db *gorm.DB, id int) *Book {
	var book Book
	result := db.First(&book, id)
	if result.Error != nil {
		log.Fatalf("Error finding book: %v", result.Error)
	}
	return &book
}

func updateBook(db *gorm.DB, book *Book) error {
	result := db.Model(&book).Updates(book)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func deleteBook(db *gorm.DB, id int) error {
	var book Book
	result := db.Delete(&book, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func deleteBook_unscope(db *gorm.DB, id uint) {
	var book Book
	result := db.Unscoped().Delete(&book, id)
	if result.Error != nil {
		log.Fatalf("Error deleting book: %v", result.Error)
	}
	fmt.Println("Book deleted successfully")
}
