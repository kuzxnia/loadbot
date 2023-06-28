package database

import "go.mongodb.org/mongo-driver/bson"

type Book struct {
	Title  string
	Author string
	ISBN   string
}

var (
	SingleBook = Book{
		Title:  "The Art of Computer Programming, Vol. 1",
		Author: "Donald E. Knuth",
		ISBN:   "978-0201896831",
	}
	MultipleBooks = []interface{}{
		Book{Title: "The Trial", Author: "Franz Kafka", ISBN: "978-0307595119"},
		Book{Title: "The Castle", Author: "Json Malone", ISBN: "978-0307474670"},
		Book{Title: "The Hunger Games", Author: "Suzanne Collins", ISBN: "978-0439023481"},
		Book{Title: "Catching Fire", Author: "Suzanne Collins", ISBN: "978-0439023498"},
		Book{Title: "A Game of Thrones", Author: "George R. R. Martin", ISBN: "978-0553593716"},
		Book{Title: "The Name of the Wind", Author: "Patrick Rothfuss", ISBN: "978-0756404741"},
		Book{Title: "Slaughterhouse-Five", Author: "Kurt Vonnegut", ISBN: "978-0385333849"},
		Book{Title: "Watchmen", Author: "Alan Moore", ISBN: "978-0930289232"},
		Book{Title: "Quicksilver", Author: "Neal Stephenson", ISBN: "978-0380977425"},
		Book{Title: "Ender's Game", Author: "Orson Scott Card", ISBN: "978-0812550702"},
	}
	FilterByTitle = bson.M{"author": "Suzanne Collins"}
)
