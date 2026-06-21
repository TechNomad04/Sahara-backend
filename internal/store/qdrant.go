package store

import "github.com/qdrant/go-client/qdrant"

type Qdrant struct {
  Client *qdrant.Client
}

func NewQdrant() (*Qdrant) {
    client, err := qdrant.NewClient(&qdrant.Config{
        Host: "localhost",
        Port: 6334,
    })

    if err != nil {
      panic(err)
    }
    return &Qdrant{Client: client}
}