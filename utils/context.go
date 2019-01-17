package utils

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// DumpOutGoingContext :nodoc:
func DumpOutGoingContext(c context.Context) string {
	md, _ := metadata.FromOutgoingContext(c)
	return Dump(md)
}

// DumpIncomingContext :nodoc:
func DumpIncomingContext(c context.Context) string {
	md, _ := metadata.FromIncomingContext(c)
	return Dump(md)
}
