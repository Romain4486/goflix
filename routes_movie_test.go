package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStore struct {
	movieId int64
	movies  []*Movie
}

func (t testStore) Open() error {
	return nil
}

func (t testStore) Close() error {
	return nil
}

func (t testStore) GetMovies() ([]*Movie, error) {
	return t.movies, nil
}

func (t testStore) GetMovieById(id int64) (*Movie, error) {
	for _, m := range t.movies {
		if m.Id == id {
			return m, nil
		}
	}
	return nil, nil
}

func (t testStore) CreateMovie(m *Movie) error {
	t.movieId++
	m.Id = t.movieId
	t.movies = append(t.movies, m)
	return nil
}

func TestMovieCreateUnit(t *testing.T) {
	srv := newServer()
	srv.store = &testStore{}
	p := struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerUrl  string `json:"trailer_url"`
	}{
		Title:       "Inception",
		ReleaseDate: "2010-07-18",
		Duration:    148,
		TrailerUrl:  "hhtp://url",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/movies/", &buf)
	w := httptest.NewRecorder()
	srv.handleMovieCreate()(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestMovieCreateIntegration(t *testing.T) {
	srv := newServer()
	srv.store = &testStore{}
	p := struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerUrl  string `json:"trailer_url"`
	}{
		Title:       "Inception",
		ReleaseDate: "2010-07-18",
		Duration:    148,
		TrailerUrl:  "hhtp://url",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(p)
	assert.Nil(t, err)

	r := httptest.NewRequest("POST", "/api/movies/", &buf)
	w := httptest.NewRecorder()
	srv.serveHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
