package graph

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

func MarshalUUID(u uuid.UUID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, fmt.Sprintf("%q", u.String()))
	})
}

func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	if s, ok := v.(string); ok {
		return uuid.Parse(s)
	}
	return uuid.Nil, fmt.Errorf("UUID must be a string")
}

func MarshalTime(t time.Time) graphql.Marshaler {
	return graphql.MarshalTime(t)
}

func UnmarshalTime(v interface{}) (time.Time, error) {
	return graphql.UnmarshalTime(v)
}

func MarshalMap(m map[string]interface{}) graphql.Marshaler {
	return graphql.MarshalMap(m)
}

func UnmarshalMap(v interface{}) (map[string]interface{}, error) {
	return graphql.UnmarshalMap(v)
}
