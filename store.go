package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	Open() error
	Close() error
	GetMovies() ([]*Movie, error)
	GetMovieById(id int64) (*Movie, error)
	CreateMovie(m *Movie) error
}

type dbStore struct {
	db *sqlx.DB
}

var schema = `CREATE TABLE IF NOT EXISTS movie(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	release_date TEXT,
	duration INTEGER,
	trailer_url TEXT
)
`

func (store *dbStore) Open() error {
	db, err := sqlx.Connect("sqlite3", "goflix.db")
	if err != nil {
		return err
	}
	log.Println("Connect to DB")
	db.MustExec(schema)
	store.db = db
	return nil
}

func (store *dbStore) Close() error {
	return store.db.Close()
}

func (store *dbStore) GetMovies() ([]*Movie, error) {
	var movies []*Movie
	err := store.db.Select(&movies, "SELECT * from movie")
	if err != nil {
		return movies, err
	}
	return movies, nil
}

func (store *dbStore) GetMovieById(id int64) (*Movie, error) {
	var m = Movie{}
	row := store.db.QueryRow("SELECT * from movie where id= $1", id)
	err := row.Scan(&m.Id, &m.Title, &m.ReleaseDate, &m.Duration, &m.TrailerUrl)
	if err != nil {
		return &m, err
	}
	return &m, nil
}

func (store *dbStore) CreateMovie(m *Movie) error {
	res, err := store.db.Exec("INSERT INTO movie (title, release_date, duration, trailer_url) VALUES (?,?,?,?)",
		m.Title, m.ReleaseDate, m.Duration, m.TrailerUrl)
	if err != nil {
		return err
	}

	m.Id, err = res.LastInsertId()
	return err
}
