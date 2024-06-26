package server

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"time"

	"github.com/anhdt1911/scraper/internal/scraper"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

func (s *Server) BatchScrape(c *gin.Context) {
	userID := c.PostForm("userID")
	f, _ := c.FormFile("file")
	file, err := f.Open()
	if err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	keywords, err := csvReader.ReadAll()
	if err != nil {
		c.JSON(500, gin.H{"msg": err})
		return
	}
	if len(keywords[0]) > 100 {
		c.JSON(400, gin.H{"msg": "over 100 keyword limits"})
		return
	}

	go func() {
		for _, v := range keywords[0] {
			if v == "" {
				continue
			}
			result, err := s.scraper.Scrape(v)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Insert to database
			_, err = s.db.Exec(c, "INSERT INTO search_result (keyword, html_content, link_amount, adword_amount, total_search_result, user_id) VALUES ($1, $2, $3, $4, $5, $6)",
				result.Keyword, result.HtmlContent, result.LinkAmount, result.AdwordAMount, result.TotalSearchResult, userID,
			)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Sleep from 1 to 10 second to hide bot scraping behavior.
			time.Sleep(time.Duration(1+rand.Intn(10)) * time.Second)
		}
	}()

	c.JSON(200, gin.H{"data": keywords})
}

func (s *Server) ScrapeResult(c *gin.Context) {
	keyword := c.PostForm("keyword")
	userID := c.PostForm("userID")
	if keyword == "" {
		c.JSON(400, gin.H{"msg": "no keyword provide"})
		return
	}
	result, err := s.scraper.Scrape(keyword)
	if err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Insert to database
	_, err = s.db.Exec(c, "INSERT INTO search_result (keyword, html_content, link_amount, adword_amount, total_search_result, user_id) VALUES ($1, $2, $3, $4, $5, $6)",
		result.Keyword, result.HtmlContent, result.LinkAmount, result.AdwordAMount, result.TotalSearchResult, userID,
	)
	if err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"data": result})
}

func (s *Server) GetSearchResultByKeyword(c *gin.Context) {
	keyID := c.Param("keyID")
	var res scraper.SearchResult
	err := s.db.QueryRow(c, "SELECT * FROM search_result WHERE id = $1", keyID).
		Scan(&res.ID, &res.Keyword, &res.HtmlContent, &res.AdwordAMount, &res.TotalSearchResult, &res.UserID, &res.LinkAmount)
	if err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	// c.JSON(200, gin.H{"data": res})
	c.Data(200, "text/html; charset=utf-8", []byte(res.HtmlContent))
}

func (s *Server) GetSearchResultsByUserID(c *gin.Context) {
	userID := c.Param("userID")
	search := c.Query("search")

	var results []scraper.SearchResult
	var rows pgx.Rows
	var err error
	if search == "" {
		rows, err = s.db.Query(c, "SELECT * FROM search_result WHERE user_id = $1 ORDER BY id DESC", userID)
	} else {
		rows, err = s.db.Query(c, "SELECT * FROM search_result WHERE user_id = $1 AND keyword LIKE '%' || $2 || '%' ORDER BY id DESC", userID, search)
	}
	if err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var res scraper.SearchResult
		if err := rows.Scan(&res.ID, &res.Keyword, &res.HtmlContent, &res.AdwordAMount, &res.TotalSearchResult, &res.UserID, &res.LinkAmount); err != nil {
			fmt.Println(err)
			return
		}
		results = append(results, res)
	}
	if err := rows.Err(); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"data": results})
}
