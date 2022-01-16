package test

import (
	"context"
	"exposures/ent"
	"testing"

	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func SetupTestWithEntClient(t *testing.T) (*ent.Client, context.Context) {
	// Create an ent.Client with in-memory SQLite database.
	client, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed opening connection to sqlite: %v", err)
	}
	ctx := context.Background()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(ctx); err != nil {
		t.Fatalf("failed creating schema resources: %v", err)
	}
	return client, ctx
}
