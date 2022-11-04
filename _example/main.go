package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	Id        int     `json:"id"`
	Title     string  `json:"title"`
	Desc      string  `json:"desc"`
	Content   string  `json:"content"`
	CreatedAt IsoTime `json:"created_at"`
}

type Response struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

const (
	INVALID_ID = "Invalid id"
	NOT_FOUND  = "Entity not found"
)

var INITIAL_ARTICLES = []Article{
	{
		Id:        0,
		Title:     "First article",
		Desc:      "The first article",
		Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		CreatedAt: IsoTimeNow(),
	},
	{
		Id:        1,
		Title:     "Second article",
		Desc:      "The second article",
		Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		CreatedAt: IsoTimeNow(),
	},
	{
		Id:        2,
		Title:     "Third article",
		Desc:      "The third article",
		Content:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		CreatedAt: IsoTimeNow(),
	},
}
var articles = INITIAL_ARTICLES

func findArticleById(id int) (int, *Article) {
	for index, article := range articles {
		if article.Id == id {
			return index, &article
		}
	}

	return -1, nil
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	w.Write([]byte("\n"))
}

func main() {
	logDest, err := os.OpenFile("access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Cannot open log file 'access.log'")
		return
	}

	middleware.DefaultLogger = middleware.RequestLogger(
		&middleware.DefaultLogFormatter{
			Logger:  log.New(logDest, "", log.LstdFlags),
			NoColor: true,
		},
	)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/articles/all", func(w http.ResponseWriter, r *http.Request) {
		respondwithJSON(w, 200, Response{
			Data: articles,
		})
	})

	r.Post("/articles", func(w http.ResponseWriter, r *http.Request) {
		article := Article{}
		json.NewDecoder(r.Body).Decode(&article)
		article.Id = len(articles) + 1
		articles = append(articles, article)

		respondwithJSON(w, 201, Response{
			Data: article,
		})
	})

	r.Get("/articles/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			respondwithJSON(w, 500, Response{
				Error: INVALID_ID,
			})
			return
		}

		_, article := findArticleById(id)
		if nil == article {
			respondwithJSON(w, 404, Response{
				Error: NOT_FOUND,
			})
			return
		}

		respondwithJSON(w, 200, Response{
			Data: article,
		})
	})

	r.Put("/articles/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			respondwithJSON(w, 500, Response{
				Error: INVALID_ID,
			})
			return
		}

		index, article := findArticleById(id)
		if index == -1 {
			respondwithJSON(w, 404, Response{
				Error: NOT_FOUND,
			})
			return
		}

		json.NewDecoder(r.Body).Decode(article)
		articles[index] = *article

		respondwithJSON(w, 200, Response{
			Data: article,
		})
	})

	r.Delete("/articles/{id:[0-9]}", func(w http.ResponseWriter, r *http.Request) {
		strId := chi.URLParam(r, "id")
		id, err := strconv.Atoi(strId)
		if err != nil {
			respondwithJSON(w, 500, Response{
				Error: INVALID_ID,
			})
			return
		}

		index, _ := findArticleById(id)
		if index == -1 {
			respondwithJSON(w, 404, Response{
				Error: NOT_FOUND,
			})
			return
		}

		articles = append(articles[:index], articles[index+1:]...)

		respondwithJSON(w, 204, nil)
	})

	r.Get("/_reset", func(w http.ResponseWriter, r *http.Request) {
		articles = make([]Article, len(INITIAL_ARTICLES))
		copy(articles, INITIAL_ARTICLES)
		respondwithJSON(w, 200, "done")
	})

	http.ListenAndServe(":3000", r)
}

type IsoTime struct {
	t     time.Time
	valid bool
}

const layout = "2006-01-02T15:04:05+0000"

func NewIsoTime(year, month, day int) IsoTime {
	if year == 0 || month == 0 {
		return IsoTime{}
	}
	if year != 0 && month != 0 && day == 0 {
		day = 10
	}
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return IsoTime{t: t, valid: true}

}

func IsoTimeNow() IsoTime {
	return IsoTime{t: time.Now(), valid: true}

}

func (mt *IsoTime) UnmarshalJSON(data []byte) error {
	stringData := strings.Trim(string(data), "\"")
	if stringData == "" {
		return nil
	}
	t, err := time.Parse(layout, stringData)
	if err != nil {
		return err
	}
	mt.t = t
	mt.valid = true
	return nil
}

func (mt IsoTime) MarshalJSON() ([]byte, error) {
	if mt.IsValid() {
		return []byte(`"` + mt.t.Format(layout) + `"`), nil
	}
	return []byte(`""`), nil
}

func (mt *IsoTime) IsValid() bool {
	return mt.valid
}
