package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type pagingResult struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	PrevPage  int   `json:"prevPage"`
	NextPage  int   `json:"nextPage"`
	Count     int64 `json:"count"`
	TotalPage int   `json:"totalPage"`
}

type pagination struct {
	ctx     *fiber.Ctx
	query   *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {
	// 1. set page and limit
	page, _ := strconv.Atoi(p.ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(p.ctx.Query("limit", "5"))

	// 2. count records
	ch := make(chan int64)
	go p.countRecords(ch)
	count := <-ch
	// var count int64
	// p.query.Model(p.records).Count(&count)

	// 3. Find Records
	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.records)

	// 4. total page (จำนวนข้อมูลทั้งหมด / จำนวนที่แสดงผลต่อpage) *ปัดเศษขึ้น
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	// 5. Find nextPage
	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}
}

func (p *pagination) countRecords(ch chan int64) {
	var count int64
	p.query.Model(p.records).Count(&count)

	ch <- count
}
