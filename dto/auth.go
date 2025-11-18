package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"requeired,email"`
	Password string `json:"password" binding:"required"`
}

type LoginWithGoogleRequest struct {
	IdToken string `json:"id_token" binding:"required"`
}

type LoginWithPasswordRequest struct {
	Email    string `json:"email" binding:"requeired,email"`
	Password string `json:"password" binding:"requeired"`
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
