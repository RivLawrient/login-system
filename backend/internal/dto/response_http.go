package dto

type ResponseWeb[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type RegisterRes struct {
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

type RefreshRes struct {
	AccessToken string `json:"access_token"`
}

type MeRes struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
