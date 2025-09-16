package helper

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func Paginator(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		pageSize, err := strconv.Atoi(ctx.Value(KeyCon("page_size")).(string))

		if err != nil || pageSize <= 0 {
			pageSize = 25
		}

		page, err := strconv.Atoi(ctx.Value(KeyCon("page")).(string))

		if err != nil || page <= 0 {
			page = 1
		}

		offset := (page - 1) * pageSize

		log.WithFields(log.Fields{
			"page":   page,
			"offset": offset,
		}).Info("set pagination")
		return db.Offset(offset).Limit(pageSize)
	}
}
