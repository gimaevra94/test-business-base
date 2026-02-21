package consts

const (
	LoginHTML           = "login.html"
	UsersSelectQuery    = "select " + UID + ", " + Name + ", " + Role + " from users"
	RequestsSelectQuery = "select " + ClientName + ", " + Phone + ", " + Address + ", " + ProblemText + " from requests"
	RequestInsertQuery  = "inset into requests (" + ClientName + ", " + Phone + ", " + Address + ", " + ProblemText + ", " + Status + ") VALUES (?, ?, ?, ?, " + "'" + New + "')"
)

const (
	Dashboard = "/dashboard"
	Login     = "/login"
	Logout    = "/logout"
	Home      = "/"
	Create    = "/create"
	Status    = "/status"
)

const (
	EmptyDB       = "The database is empty"
	InternalError = "Internal server error"
	BadRequest    = "Bad request"
	BadInput      = "Bad input"
	NotAllowed    = "Method not allowed"
)

const (
	UID         = "user_id"
	Role        = "role"
	Name        = "name"
	ClientName  = "client_name"
	Phone       = "phone"
	Address     = "address"
	ProblemText = "problem_text"
	New         = "new"
)
