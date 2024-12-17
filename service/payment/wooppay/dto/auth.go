package dto

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type AuthErrorResponse struct {
	Error AuthError `json:"error"`
}

type AuthError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
