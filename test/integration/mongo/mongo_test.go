package mongo

import (
	"github.com/danielpacak/myevents-events-service/config"
	"github.com/danielpacak/myevents-events-service/domain"
	"github.com/danielpacak/myevents-events-service/persistence/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestEventsRepository(t *testing.T) {
	// FIXME Do not skip this test
	t.Skip("Configure Docker Compose or use docker-skd")
	repo, err := mongo.NewMongoEventsRepository(&config.MongoDBConfig{
		ConnectionURL: "mongodb://127.0.0.1",
		DatabaseName:  "testdb",
	})
	require.NoError(t, err)

	events, err := repo.FindAll()
	assert.Empty(t, events)

	eventId, err := repo.Create(domain.Event{
		Name:      "Czarownice z Eastwick",
		StartDate: 1234243,
		EndDate:   341342,
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

	foundById, err := repo.FindById(eventId)
	require.NoError(t, err)

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
