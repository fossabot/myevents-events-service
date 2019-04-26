package mongo

import (
	"encoding/hex"
	"github.com/danielpacak/myevents-events-service/pkg/config"
	"github.com/danielpacak/myevents-events-service/pkg/domain"
	"github.com/danielpacak/myevents-events-service/pkg/persistence/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestEventsRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test")
	}
	repo, err := mongo.NewMongoEventsRepository(config.MongoDBConfig{
		ConnectionURL: "mongodb://mongo:27017",
		DatabaseName:  "testdb",
	})
	require.NoError(t, err)

	events, err := repo.FindAll()
	assert.Empty(t, events)

	t.Logf("Creating event")
	eventId, err := repo.Create(domain.Event{
		Name:      "Czarownice z Eastwick",
		StartDate: time.Now(),
		EndDate:   time.Now(),
		Duration:  34234,
		Location: domain.Location{
			Name:      "Teatr Syrena",
			Address:   "ul. Litewska 3, 00-589 Warszawa",
			Country:   "Polska",
			OpenTime:  23434,
			CloseTime: 342343,
			Halls: []domain.Hall{
				{
					Name:     "Sala Główna",
					Location: "Parter",
				},
			},
		},
	})
	require.NoError(t, err)
	t.Logf("Done creating event %s", hex.EncodeToString(eventId))

	t.Logf("Finding event by ID: %s", hex.EncodeToString(eventId))
	foundById, err := repo.FindById(eventId)
	require.NoError(t, err)
	t.Logf("Done finding event by ID: %s", hex.EncodeToString(eventId))

	assert.Equal(t, bson.ObjectId(eventId), foundById.ID)
	assert.Equal(t, "Czarownice z Eastwick", foundById.Name)
	assert.NotNil(t, foundById.Location.ID)
	assert.Equal(t, "Teatr Syrena", foundById.Location.Name)

	foundByName, err := repo.FindByName("Czarownice z Eastwick")
	require.NoError(t, err)
	assert.Equal(t, bson.ObjectId(eventId), foundByName.ID)
	assert.Equal(t, "Czarownice z Eastwick", foundByName.Name)

	eventId, err = repo.Create(domain.Event{
		Name: "Ferdydurke",
		Location: domain.Location{
			Name:    "Scena na Woli",
			Address: "Kasprzaka 22, 01-211 Warszawa",
			Country: "Polska",
		},
	})
	require.NoError(t, err)

	events, err = repo.FindAll()
	require.NoError(t, err)
	assert.Equal(t, 2, len(events))
}
