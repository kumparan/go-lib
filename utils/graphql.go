package utils

import (
	"fmt"

	graphql "github.com/graph-gophers/graphql-go"
)

// Int642GraphQLID :nodoc:
func Int642GraphQLID(i int64) graphql.ID {
	return graphql.ID(fmt.Sprintf("%d", i))
}
