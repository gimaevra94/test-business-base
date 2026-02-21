package structs

type Request struct {
	ID, AssignedTo, Version                                               int
	ClientName, Phone, Address, ProblemText, Status, CreatedAt, UpdatedAt string
}

type User struct {
	ID   int
	Name string
	Role string
}

type LoginData struct {
	Users []User
	Msg string
}
