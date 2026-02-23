package structs

import "database/sql"

type Request struct {
	ID, Version                                                           int
	ClientName, Phone, Address, ProblemText, Status, CreatedAt, UpdatedAt string
	AssignedTo                                                            sql.NullInt64
}

type User struct {
	ID   int
	Name string
	Role string
}

type LoginData struct {
	Users    []User
	Msg      string
	User     string
	Role     string
	Requests []Request
	Masters  []User
	Error    string
}
