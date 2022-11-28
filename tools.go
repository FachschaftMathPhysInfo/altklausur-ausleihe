//go:build tools
// +build tools

package tools

import (
	// used to generate the graphql code in the backend
	_ "github.com/99designs/gqlgen"
	_ "github.com/99designs/gqlgen/graphql/introspection"
)
