package gormdb

import "github.com/alioygur/gocart/engine"
import "github.com/jinzhu/gorm"

func translateFilter(db *gorm.DB, f []*engine.Filter) *gorm.DB {
	for _, f := range f {
		switch f.Condition {
		case engine.Equal:
			db = db.Where(f.Property+" = ?", f.Value)
		case engine.In:
			db = db.Where(f.Property+" in (?)", f.Value)
		case engine.LessThan:
			db = db.Where(f.Property+" < ?", f.Value)
		case engine.LessThanOrEqual:
			db = db.Where(f.Property+" <= ?", f.Value)
		case engine.GreaterThan:
			db = db.Where(f.Property+" > ?", f.Value)
		case engine.GreaterThanOrEqual:
			db = db.Where(f.Property+" >= ?", f.Value)
		}
	}
	return db
}

func translateQuery(db *gorm.DB, q *engine.Query) *gorm.DB {
	db = translateFilter(db, q.Filters)

	for _, order := range q.Orders {
		switch order.Direction {
		case engine.Ascending:
			db = db.Order(order.Property)
		case engine.Descending:
			db = db.Order(order.Property + " DESC")
		}
	}

	if q.Limit != nil {
		db = db.Offset(q.Limit.Offset)
		db = db.Limit(q.Limit.Limit)
	}

	return db
}
