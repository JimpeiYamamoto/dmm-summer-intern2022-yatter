package utils

import (
	"strconv"
	"yatter-backend-go/app/domain/object"
)

func FixQuery(q object.Query) (object.Query, error) {
	if q.MaxID == "" {
		q.MaxID = "10000000"
	}
	if q.SinceID == "" {
		q.SinceID = "0"
	}
	if q.Limit == "" {
		q.Limit = "40"
	} else {
		limit, err := strconv.ParseUint(q.Limit, 10, 64)
		if err != nil {
			return q, err
		}
		if limit > 80 {
			q.Limit = "80"
		}
	}
	return q, nil
}
