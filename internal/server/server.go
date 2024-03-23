package server

import (
	"github.com/anhdt1911/scraper/internal/scraper"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db      *pgxpool.Pool
	scraper *scraper.Scraper
}

func New(db *pgxpool.Pool, scraper *scraper.Scraper) *Server {
	return &Server{
		db,
		scraper,
	}
}

func GetSearchResultByKeyword(ctx *gin.Context) {}
