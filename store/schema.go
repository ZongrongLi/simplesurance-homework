package store

import (
	"time"
)

type CountStatistic struct {
	CreatedAt time.Time
	Count     int64
}
