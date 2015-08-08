package search

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoSearcher struct {
	db *mgo.Database
}

func NewMongoSearcher(db *mgo.Database) *MongoSearcher {
	return &MongoSearcher{db}
}

func (m *MongoSearcher) CreateQuery(query Query) *mgo.Query {
	c := m.db.C(MongoCollectionNames[query.Resource])
	return c.Find(m.CreateQueryObject(query))
}

func (m *MongoSearcher) CreateQueryObject(query Query) bson.M {
	var result bson.M
	for _, p := range query.Params() {
		switch p := p.(type) {
		//case CompositeParam:
		//
		case *DateParam:
			//
		case *NumberParam:
			//
		case *QuantityParam:
			//
		case *ReferenceParam:
			result = m.CreateReferenceQueryObject(p)
		case *StringParam:
			result = m.CreateStringQueryObject(p)
		case *TokenParam:
			result = m.CreateTokenQueryObject(p)
		case *URIParam:
			//
		default:
			result = bson.M{}
		}
	}

	return result
}

func (m *MongoSearcher) CreateReferenceQueryObject(r *ReferenceParam) bson.M {
	single := func(path string, dtype string) bson.M {
		result := bson.M{}
		id, ok := r.GetId()
		if ok {
			result[fmt.Sprintf("%s.referenceid", path)] = ci(id)
			typ, ok := r.GetType()
			if ok {
				result[fmt.Sprintf("%s.type", path)] = typ
			}
		} else {
			url, ok := r.GetUrl()
			if ok {
				result[fmt.Sprintf("%s.reference", path)] = ci(url)
			}
		}
		return result
	}
	return orPaths(single, r.Paths)
}

func (m *MongoSearcher) CreateStringQueryObject(s *StringParam) bson.M {
	single := func(path string, dtype string) bson.M {
		result := bson.M{}
		if s.Name == "_id" {
			result[path] = ci(s.String)
		} else {
			result[path] = cisw(s.String)
		}
		return result
	}

	return orPaths(single, s.Paths)
}

func (m *MongoSearcher) CreateTokenQueryObject(t *TokenParam) bson.M {
	single := func(path string, dtype string) bson.M {
		result := bson.M{}
		switch dtype {
		case "Coding":
			result[fmt.Sprintf("%s.code", path)] = ci(t.Code)
			if !t.AnySystem {
				result[fmt.Sprintf("%s.system", path)] = ci(t.System)
			}
		case "CodeableConcept":
			if t.AnySystem {
				result[fmt.Sprintf("%s.coding.code", path)] = ci(t.Code)
			} else {
				result[fmt.Sprintf("%s.coding", path)] = bson.M{"$elemMatch": bson.M{"system": ci(t.System), "code": ci(t.Code)}}
			}
		case "Identifier":
			result[fmt.Sprintf("%s.value", path)] = ci(t.Code)
			if !t.AnySystem {
				result[fmt.Sprintf("%s.system", path)] = ci(t.System)
			}
		case "ContactPoint":
			result[fmt.Sprintf("%s.value", path)] = ci(t.Code)
			if !t.AnySystem {
				result[fmt.Sprintf("%s.use", path)] = ci(t.System)
			}
		case "code", "boolean", "string":
			result[path] = ci(t.Code)
		}

		return result
	}

	return orPaths(single, t.Paths)
}

// When multiple paths are present, they should be represented as an OR.
// objFunc is a function that generates a single query for a path
func orPaths(objFunc func(string, string) bson.M, paths map[string]string) bson.M {
	results := make([]bson.M, 0, len(paths))
	for k, v := range paths {
		results = append(results, objFunc(k, v))
	}

	if len(results) == 1 {
		return results[0]
	} else {
		return bson.M{"$or": results}
	}
}

// Case-insensitive match
func ci(s string) bson.RegEx {
	return bson.RegEx{Pattern: fmt.Sprintf("^%s$", s), Options: "i"}
}

// Case-insensitive starts-with
func cisw(s string) bson.RegEx {
	return bson.RegEx{Pattern: fmt.Sprintf("^%s", s), Options: "i"}
}
