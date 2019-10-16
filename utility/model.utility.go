package utility

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// MongoSearchFieldParser ...
func MongoSearchFieldParser(fields []string, keyword string) []bson.M {
	var or []bson.M

	for _, f := range fields {
		or = append(or, bson.M{f: bson.M{"$regex": keyword}})
	}

	return or
}

// GetOrderByKeyAndValue ...
func GetOrderByKeyAndValue(orderBy string) (string, int) {
	if orderBy != "" {
		var orderKey string
		var orderValue int

		if strings.Contains(orderBy, "_ASC") {
			s := strings.Split(orderBy, "_ASC")
			orderKey = s[0]
			orderValue = 1
		} else if strings.Contains(orderBy, "_DESC") {
			s := strings.Split(orderBy, "_DESC")
			orderKey = s[0]
			orderValue = -1
		}

		return orderKey, orderValue
	}
	return "created_at", 1
}
