package store

import (
	"context"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

type Qdrant struct {
  Client *qdrant.Client
}

func NewQdrant() (*Qdrant) {
    client, err := qdrant.NewClient(&qdrant.Config{
        Host: "localhost",
        Port: 6334,
    })

    if err != nil {
      log.Fatal(err)
    }
    return &Qdrant{Client: client}
}

func CreateCollection(q *Qdrant) {
  err := q.Client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "documents",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     384,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		log.Fatal(err)
	}
}