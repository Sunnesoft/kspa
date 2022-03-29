package graph_shortest_paths

import (
	"math/rand"
	"testing"
	"time"
)

func setupTestCase(t *testing.T) func(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	return func(t *testing.T) {
	}
}
