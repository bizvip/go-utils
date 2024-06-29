/******************************************************************************
 * Copyright (c) 2024. Archer++. All rights reserved.                         *
 * Author ORCID: https://orcid.org/0009-0003-8150-367X                        *
 ******************************************************************************/

package repo

import (
	"time"

	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
)

type User struct {
	ID       int
	JoinedAt time.Time `gorm:"column:created_at"`
}

func CreateUserPaginator(cursor paginator.Cursor, order *paginator.Order, limit *int) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			Keys:  []string{"ID", "JoinedAt"},
			Limit: 10,
			Order: paginator.ASC,
		},
	}
	if limit != nil {
		opts = append(opts, paginator.WithLimit(*limit))
	}
	if order != nil {
		opts = append(opts, paginator.WithOrder(*order))
	}
	if cursor.After != nil {
		opts = append(opts, paginator.WithAfter(*cursor.After))
	}
	if cursor.Before != nil {
		opts = append(opts, paginator.WithBefore(*cursor.Before))
	}
	return paginator.New(opts...)
}

// func CreateUserPaginator(/* ... */) {
//	p := paginator.New(
//		&paginator.Config{
//			Rules: []paginator.Rule{
//				{
//					Key: "ID",
//				},
//				{
//					Key: "JoinedAt",
//					Order: paginator.DESC,
//					SQLRepr: "users.created_at",
//					NULLReplacement: "1970-01-01",
//				},
//			},
//			Limit: 10,
//			// Order here will apply to keys without order specified.
//			// In this example paginator will order by "ID" ASC, "JoinedAt" DESC.
//			Order: paginator.ASC,
//		},
//	)
//	// ...
//	return p
// }

// func CreateUserPaginator(
//
//	cursor paginator.Cursor,
//	order *paginator.Order,
//	limit *int,
//
//	) *paginator.Paginator {
//		p := paginator.New(
//			&paginator.Config{
//				Keys: []string{"ID", "JoinedAt"},
//				Limit: 10,
//				Order: paginator.ASC,
//			},
//		)
//		if order != nil {
//			p.SetOrder(*order)
//		}
//		if limit != nil {
//			p.SetLimit(*limit)
//		}
//		if cursor.After != nil {
//			p.SetAfterCursor(*cursor.After)
//		}
//		if cursor.Before != nil {
//			p.SetBeforeCursor(*cursor.Before)
//		}
//		return p
//	}

type Cursor struct {
	After  *string `json:"after" query:"after"`
	Before *string `json:"before" query:"before"`
}
