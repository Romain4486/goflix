package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type jsonMovie struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration    int    `json:"duration"`
	TrailerUrl  string `json:"trailer_url"`
}

func (s *server) handleMovieList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := s.store.GetMovies()
		if err != nil {
			log.Printf("Cannot load movies. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		var resp = make([]jsonMovie, len(movies))
		for i, m := range movies {
			resp[i] = mapMovieToJson(m)
		}
		s.respond(w, r, resp, http.StatusOK)
	}
}

func (s *server) handleMovieById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		id, err := strconv.ParseInt(params["id"], 10, 64)
		log.Printf("id=%v", id)
		if err != nil {
			log.Printf("The id specified isn't an integer err=%v\n", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		m, err := s.store.GetMovieById(id)
		if err != nil {
			log.Printf("Cannot load movie. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		jsonMovie := mapMovieToJson(m)
		s.respond(w, r, jsonMovie, http.StatusOK)
	}
}

func (s *server) handleMovieCreate() http.HandlerFunc {
	type request struct {
		Title       string `json:"title"`
		ReleaseDate string `json:"release_date"`
		Duration    int    `json:"duration"`
		TrailerUrl  string `json:"trailer_url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		err := s.decode(w, r, &req)
		if err != nil {
			log.Printf("Cannot parse json to movie. err=%v", err)
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}
		m := &Movie{
			Id:          0,
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
			Duration:    req.Duration,
			TrailerUrl:  req.TrailerUrl,
		}

		err = s.store.CreateMovie(m)
		if err != nil {
			log.Printf("Cannot create movie in DB. err=%v\n", err)
			s.respond(w, r, nil, http.StatusInternalServerError)
			return
		}

		jsonMovie := mapMovieToJson(m)
		s.respond(w, r, jsonMovie, http.StatusOK)
	}
}

func mapMovieToJson(m *Movie) jsonMovie {
	return jsonMovie{
		Id:          m.Id,
		Title:       m.Title,
		ReleaseDate: m.ReleaseDate,
		Duration:    m.Duration,
		TrailerUrl:  m.TrailerUrl,
	}
}
