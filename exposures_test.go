// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"exposures/ent"
	"exposures/messages"
	"exposures/requesthandlers"
	"exposures/smshandler"
	"net/http"
	"net/http/httptest"
	"testing"

	"entgo.io/ent/dialect"
)

func setupTestWithEntClient(t *testing.T) (*ent.Client, context.Context) {
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

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(requesthandlers.IndexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusOK,
		)
	}
}

func TestIndexHandlerNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(requesthandlers.IndexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf(
			"unexpected status: got (%v) want (%v)",
			status,
			http.StatusNotFound,
		)
	}
}

func TestRsvp(t *testing.T) {
	client, ctx := setupTestWithEntClient(t)
	defer client.Close()

	sampleCode := "205202ac-c4e6-48ed-b469-9b3bcf592316"
	ms := smshandler.Rsvp(ctx, client, "+1(123)456-7890", sampleCode)
	if len(ms) != 1 {
		t.Errorf(
			"Expected RSVP to send 1 message. Sent %v.", ms,
		)
	}

	if ms[0].Type != messages.RsvpSuccess {
		t.Errorf(
			"Expected RSVP message, but actually got: %v", ms,
		)
	}

	users := client.User.Query().AllX(ctx)

	if len(users) != 1 {
		t.Errorf(
			"Expected 1 user to be created and persisted to DB, but actually got: %v", len(users),
		)
	}
}

func TestHandlePositiveCase(t *testing.T) {
	client, ctx := setupTestWithEntClient(t)
	defer client.Close()

	user1 := "+1(123)456-7890"
	user2 := "+1(234)567-8901"
	user3 := "+1(345)678-9012"
	sampleCode := "205202ac-c4e6-48ed-b469-9b3bcf592316"
	smshandler.Rsvp(ctx, client, user1, sampleCode)
	smshandler.Rsvp(ctx, client, user2, sampleCode)
	smshandler.Rsvp(ctx, client, user3, sampleCode)

	ms := smshandler.HandlePositiveCase(ctx, client, user1)
	if len(ms) != 3 {
		t.Fatalf("Expected 3 messages to be sent. Got: %v", len(ms))
	}

	for _, m := range ms {
		if m.Recipient == user1 && m.Type != messages.ThankForSelfReporting {
			t.Fatalf("Expected to thank user 1 for reporting. Found this message instead: %v", ms)
		}

		if m.Recipient == user2 && m.Type != messages.NotifyPositiveCase {
			t.Fatalf("Expected to notify user 2. Sent this message instead: %v", ms)
		}

		if m.Recipient == user3 && m.Type != messages.NotifyPositiveCase {
			t.Fatalf("Expected to notify user 3. Sent this message instead: %v", ms)
		}
	}
}
