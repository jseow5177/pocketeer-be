package firestore

import (
	"cloud.google.com/go/firestore"
)

type CategoryFirestore struct {
	client *firestore.Client
}

func NewCategoryFirestore(client *firestore.Client) *CategoryFirestore {
	return &CategoryFirestore{
		client: client,
	}
}
