package models

type ResponseFriend struct {
	Success bool `json:"success"`
	Friends []string `json:"friends"`
	Count int64 `json:"count"`
	Message string `json:"message"`
}

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type SubscribeResponse struct {
	Success bool `json:"success"`
	Recipients []string `json:"recipients"`
	Message string `json:"message"`
}

