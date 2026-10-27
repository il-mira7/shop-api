package handler

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email" example:"ivanov@mail.ru"`
	Name        string `json:"name" binding:"required,min=2,max=100" example:"Иван"`
	Surname     string `json:"surname" binding:"required,min=2,max=100" example:"Иванов"`
	PhoneNumber string `json:"phone_number" binding:"required,min=5" example:"+79990001122"`
	Password    string `json:"password" binding:"required,min=8" example:"secret123"`
}

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email" example:"ivanov@mail.ru"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

type ResetRequest struct {
	Email string `json:"email" binding:"required,email" example:"ivanov@mail.ru"`
}
