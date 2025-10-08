package dtos

type UserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ErrorDTO struct {
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type TodoListDTO struct {
	Name    string `json:"name"`
	ColorID int    `json:"color_id"`
}

type TodoDTO struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}
