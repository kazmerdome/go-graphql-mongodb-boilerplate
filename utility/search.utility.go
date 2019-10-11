package utility

import "go.mongodb.org/mongo-driver/bson"

// SearchParser ...
func SearchParser(fields []string, keyword string) []bson.M {
	var or []bson.M

	for _, f := range fields {
		or = append(or, bson.M{f: bson.M{"$regex": keyword}})
	}

	return or
}
