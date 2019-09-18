package utils

import (
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
)

// Int642GraphQLID :nodoc:
func Int642GraphQLID(i int64) graphql.ID {
	return graphql.ID(Int642String(i))
}

// GraphQLID2String :nodoc:
func GraphQLID2String(id graphql.ID) string {
	return string(id)
}

// GraphQLIDPointer2String :nodoc:
func GraphQLIDPointer2String(id *graphql.ID) string {
	if id == nil {
		return ""
	}

	return string(*id)
}

// GraphQLID2Int64 :nodoc:
func GraphQLID2Int64(id graphql.ID) int64 {
	return String2Int64(string(id))
}

// GraphQLIDPointer2Int64 :nodoc:
func GraphQLIDPointer2Int64(id *graphql.ID) int64 {
	if id == nil {
		return int64(0)
	}

	return String2Int64(string(*id))
}

// GraphQLID2Int32 :nodoc:
func GraphQLID2Int32(id graphql.ID) int32 {
	newID, err := strconv.Atoi(string(id))
	if err != nil {
		return int32(0)
	}
	return int32(newID)
}

// GraphQLIDPointer2Int32 :nodoc:
func GraphQLIDPointer2Int32(id *graphql.ID) int32 {
	if id == nil {
		return int32(0)
	}

	newID, err := strconv.Atoi(string(*id))
	if err != nil {
		return int32(0)
	}
	return int32(newID)
}
