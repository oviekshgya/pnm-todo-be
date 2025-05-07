package pkg

type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type CRUDProduct struct {
	ID     uint   `form:"id" json:"id"`
	Name   string `form:"name" json:"name"`
	Jumlah int    `form:"jumlah" json:"jumlah"`
}
