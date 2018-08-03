package pagination

import (
	"math"

	"github.com/jinzhu/gorm"
)

// Param struct
type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

// Pagination struct
type Pagination struct {
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Records     interface{} `json:"records"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

// Pagging func
func Pagging(p *Param, dataSource interface{}) *Pagination {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var pagination Pagination
	var count int
	var offset int

	go countRecords(db, dataSource, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(dataSource)
	<-done

	pagination.TotalRecord = count
	pagination.Records = dataSource
	pagination.Page = p.Page

	pagination.Offset = offset
	pagination.Limit = p.Limit
	pagination.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		pagination.PrevPage = p.Page - 1
	} else {
		pagination.PrevPage = p.Page
	}

	if p.Page == pagination.TotalPage {
		pagination.NextPage = p.Page
	} else {
		pagination.NextPage = p.Page + 1
	}
	return &pagination
}

func countRecords(db *gorm.DB, countDataSource interface{}, done chan bool, count *int) {
	db.Model(countDataSource).Count(count)
	done <- true
}
