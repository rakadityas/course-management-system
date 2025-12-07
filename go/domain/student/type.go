package studentdomain

import "time"

type Student struct {
	ID         int64
	Email      string
	CreateTime time.Time
	UpdateTime time.Time
}
