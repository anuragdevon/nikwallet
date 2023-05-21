package dto

type UserSignupRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupResponseDTO struct {
	UserID  int    `json:"user_id"`
	IDToken string `json:"id_token"`
}

type UserSigninRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSigninResponseDTO struct {
	IDToken string `json:"id_token"`
}
