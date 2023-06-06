package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type Book struct {
	gorm.Model
	Author    string
	Title     string
	Publisher string
	GenreID   string
	Genre     Genre `gorm:"foreignKey:GenreID"`
}

//type Order struct {
//	ID         uint
//	CustomerID uint `gorm:"foreignKey:CustomerID"` // Foreign key field
//	// Other fields
//}

type Author struct {
	gorm.Model
	AuthorID string
	Name     string
}

type Publisher struct {
	gorm.Model
	PublisherID   string
	PublisherName string
	Country       string
}

type Genre struct {
	gorm.Model
	GenreID   string
	GenreName string
}

func setup(db *gorm.DB) {
	db.AutoMigrate(&Book{}, &Author{}, &Publisher{}, &Genre{})
	seed(db)
}

func seed(db *gorm.DB) {
	var authors = []*Author{
		{AuthorID: "A1", Name: "J. K. Rowling"},
		{AuthorID: "A2", Name: "Stan Lee"},
		{AuthorID: "A3", Name: "Paula Hawkins"},
	}

	//db.Create(&authors)
	for _, author := range authors {
		db.Create(&author)
	}

	var publishers = []*Publisher{
		{PublisherID: "P1", PublisherName: "Bloomsbury Publishing", Country: "London"},
		{PublisherID: "P2", PublisherName: "Marvel Comics", Country: "United States"},
		{PublisherID: "P3", PublisherName: "Transworld Publishers Ltd", Country: "London"},
	}

	for _, publisher := range publishers {
		db.Create(&publisher)
	}

	var genres = []*Genre{
		{GenreID: "G1", GenreName: "Fiction"},
		{GenreID: "G2", GenreName: "Comic"},
		{GenreID: "G3", GenreName: "Thriller"},
	}

	for _, genre := range genres {
		db.Create(&genre)
	}

	var genreType1, genreType2, genreType3 Genre
	var publisherCountry Publisher
	db.First(&genreType1, "genre_name = ?", "Fiction")
	fmt.Println("Genre ID:", genreType1.GenreID)
	fmt.Println("Genre Name:", genreType1.GenreName)

	db.First(&genreType2, "genre_name = ?", "Comic")
	fmt.Println("Genre ID:", genreType2.GenreID)
	fmt.Println("Genre Name:", genreType2.GenreName)

	db.First(&genreType3, "genre_name = ?", "Thriller")
	fmt.Println("Genre ID:", genreType3.GenreID)
	fmt.Println("Genre Name:", genreType3.GenreName)

	db.First(&publisherCountry, "country = ?", "London")
	fmt.Println("Publisher ID:", publisherCountry.PublisherID)
	fmt.Println("Publisher Name:", publisherCountry.PublisherName)
	fmt.Println("Publisher Country:", publisherCountry.Country)

	var authorName Author
	db.First(&authorName, "name = ?", "Paula Hawkins")

	//var books = []*Book{
	//	{Title: "Harry Potter", Author: "J. K. Rowling", Publisher: "Bloomsbury Publishing", Genre: genreType.GenreName},
	//	{Title: "Spider-Man", Author: "Stan Lee", Publisher: "Marvel Comics", Genre: "Comic"},
	//	{Title: "The Girl on the Train", Author: authorName.Name, Publisher: "Transworld Publishers Ltd", Genre: "Thriller"},
	//}

	var books = []*Book{
		{Title: "Harry Potter", Author: "J. K. Rowling", Publisher: "Bloomsbury Publishing", GenreID: genreType1.GenreID, Genre: genreType1},
		//{Title: "Spider-Man", Author: "Stan Lee", Publisher: "Marvel Comics", GenreID: genreType2},
		//{Title: "The Girl on the Train", Author: authorName.Name, Publisher: "Transworld Publishers Ltd", GenreID: genreType3.GenreID},
	}

	for _, book := range books {
		db.Create(&book)
	}
}

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("Cannot connect the Database")
	}

	defer db.Close()
	db.LogMode(true)
	setup(db)

	var books []Book
	db.Model(&Book{}).Preload("Genre").Find(&books)

	for _, book := range books {
		fmt.Println("\nID: ", book.ID, "\nTitle:", book.Title, "\nAuthor:", book.Author, "\nPublisher:", book.Publisher, "\nGenre Name:", book.Genre)
	}
	doError(db)
}

func doError(db *gorm.DB) {
	var Othello Book
	if err := db.Where("Title = ?", "Othello").First(&Othello).Error; err != nil {
		log.Fatalf("Error while loading Book! %s", err)
	}
}
