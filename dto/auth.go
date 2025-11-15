package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"requeired,email"`
	Password string `json:"password" validate:"required"`
}

type LoginWithGoogleRequest struct {
	IdToken string `json:"id_token" validate:"required"`
}

type LoginWithPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Uid   string `json:"uid"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	Uid   string `json:"uid"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
