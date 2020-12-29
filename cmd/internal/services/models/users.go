package models

type Users struct {
	ID int
	Name string
	Email string
}

type Friend struct {
	UserEmail string
	FriendEmail string
}

type UserRelationship struct {
	RelationshipId int `json:"relationshipId"`
	UserEmail string `json:"userEmail"`
	FriendEmail string `json:"friendEmail"`
	IsFriend bool `json:"isFriend"`
	StatusUpdate string `json:"statusUpdate"`
}


