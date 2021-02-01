package models

type ResponseFriend struct {
	Success bool `json:"success"`
	Friends []string `json:"friends"`
	Count int64 `json:"count"`
}

type Response struct {
	Success bool `json:"success"`
}

type SubscribeResponse struct {
	Success bool `json:"success"`
	Recipients []string `json:"recipients"`
}

