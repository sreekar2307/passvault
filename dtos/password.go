package dtos

type GeneratePasswordParams struct {
	Size             int  `json:"size" binding:"required,min=8,max=64"`
	ExcludeDigits    bool `json:"excludeDigits"`
	ExcludeAlphabets bool `json:"excludeAlphabets"`
	ExcludeSymbols   bool `json:"excludeSymbols"`
}

type GetPasswordFilter struct {
	UserID      uint
	ID          uint
	NameLike    string
	WebsiteLike string
	Email       string
	Limit       int
	Offset      int
}

type GetPasswordsParams struct {
	Query  string `form:"query" binding:"omitempty,max=30"`
	Limit  int    `form:"limit" binding:"omitempty,max=20,min=0"`
	Offset int    `form:"offset" binding:"omitempty,min=0"`
}

type StorePasswordParams struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Notes    string `json:"notes"`
	Website  string `json:"website" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type UpdatePasswordParams struct {
	Name     string `json:"name" binding:"omitempty,min=1"`
	ID       uint   `uri:"id" binding:"required"`
	Username string `json:"username"`
	Password string `json:"password" binding:"omitempty,min=1"`
	Notes    string `json:"notes"`
	Website  string `json:"website" binding:"omitempty,min=1"`
	Email    string `json:"email" binding:"omitempty,min=1"`
}
