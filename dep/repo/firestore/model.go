package firestore

type Category struct {
	CatName    *string `firestore:"cat_name"`
	CreateTime *uint64 `firestore:"create_time"`
	UpdateTime *uint64 `firestore:"update_time"`
}
