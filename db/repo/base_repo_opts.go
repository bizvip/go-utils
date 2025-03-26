package repo

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// SelOpt 查询可选选项
type SelOpt func(*gorm.DB) *gorm.DB

func WithOrderBy(orderBy string) SelOpt {
	return func(q *gorm.DB) *gorm.DB { return q.Order(orderBy) }
}

func WithLimit(limit int) SelOpt {
	return func(q *gorm.DB) *gorm.DB { return q.Limit(limit) }
}

func WithWhere(conditions map[string]interface{}) SelOpt {
	return func(q *gorm.DB) *gorm.DB {
		for key, value := range conditions {
			var field, operator string

			// 分离字段名和操作符
			if parts := strings.SplitN(key, " ", 2); len(parts) == 2 {
				field = parts[0]
				operator = parts[1]
			} else {
				field = key
				operator = "=" // 默认为等于操作
			}

			// 根据操作符选择适当的查询
			switch operator {
			case ">=":
				q = q.Where(fmt.Sprintf("%s >= ?", field), value)
			case "<=":
				q = q.Where(fmt.Sprintf("%s <= ?", field), value)
			case ">":
				q = q.Where(fmt.Sprintf("%s > ?", field), value)
			case "<":
				q = q.Where(fmt.Sprintf("%s < ?", field), value)
			case "!=":
				q = q.Where(fmt.Sprintf("%s != ?", field), value)
			case "=":
				fallthrough // 直接到下一个case "="，即默认情况
			default:
				q = q.Where(fmt.Sprintf("%s = ?", field), value)
			}
		}
		return q
	}
}
func WithSelect(columns []string) SelOpt {
	return func(q *gorm.DB) *gorm.DB {
		return q.Select(columns)
	}
}
