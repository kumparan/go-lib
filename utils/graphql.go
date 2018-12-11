package utils

import (
	"fmt"
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"
)

// Int642GraphQLID :nodoc:
func Int642GraphQLID(i int64) graphql.ID {
	return graphql.ID(fmt.Sprintf("%d", i))
}

// GraphQLID2String :nodoc:
func GraphQLID2String(id graphql.ID) string {
	return fmt.Sprintf("%s", id)
}

// GraphQLIDPointer2String :nodoc:
func GraphQLIDPointer2String(id *graphql.ID) string {
	if id == nil {
		return ""
	}

	return fmt.Sprintf("%s", *id)
}

// GraphQLID2Int64 :nodoc:
func GraphQLID2Int64(id graphql.ID) int64 {
	newID, err := strconv.Atoi(fmt.Sprintf("%s", id))
	if err != nil {
		return int64(0)
	}
	return int64(newID)
}

// GraphQLIDPointer2Int64 :nodoc:
func GraphQLIDPointer2Int64(id *graphql.ID) int64 {
	if id == nil {
		return int64(0)
	}

	newID, err := strconv.Atoi(fmt.Sprintf("%s", *id))
	if err != nil {
		return int64(0)
	}
	return int64(newID)
}

// GraphQLID2Int32 :nodoc:
func GraphQLID2Int32(id graphql.ID) int32 {
	newID, err := strconv.Atoi(fmt.Sprintf("%s", id))
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

	newID, err := strconv.Atoi(fmt.Sprintf("%s", *id))
	if err != nil {
		return int32(0)
	}
	return int32(newID)
}
