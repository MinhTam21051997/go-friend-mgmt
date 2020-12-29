package models

type FriendsList struct{
	Friends []string `json:"friends"`
}

type EmailUser struct{
	Email string `json:"email"`
}

type SubscribeUser struct {
	Requestor string `json:"requestor"`
	Target string `json:"target"`
}

type RetrieveUpdate struct {
	Sender string `json:"sender"`
	Text string `json:"text"`
}
