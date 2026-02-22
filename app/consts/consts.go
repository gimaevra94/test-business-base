package consts

const (
	UsersSelectQuery            = "select " + ID + ", " + Name + ", " + Role + " from users"
	RequestsSelectQuery         = "select " + ClientName + ", " + Phone + ", " + Address + ", " + ProblemText + " from requests"
	RequestInsertQuery          = "insert into requests (" + ClientName + ", " + Phone + ", " + Address + ", " + ProblemText + ", " + Status + ") values (?, ?, ?, ?, " + "'" + New + "')"
	DashboardSelectQuery        = "select " + ID + ", " + ClientName + ", " + Phone + ", " + Address + ", " + ProblemText + ", " + Status + ", " + AssignedTo + ", " + Version + ", " + CreatedAt + ", " + UpdatedAt + " from requests WHERE 1=1"
	MastersSelectQuery          = "select " + ID + ", " + Name + " from users where role = '" + Master + "'"
	AssignedStatusUpdateQuery   = "update requests set status = '" + Assigned + "', " + AssignedTo + " = ? WHERE " + ID + " = ? AND status = '" + New + "'"
	CanceledStatusUpdateQuery   = "update requests set status = '" + Canceled + "' where " + ID + " = ? and status in ('" + New + "', '" + Assigned + "')"
	InProgressStatusUpdateQuery = "update requests set status = '" + InProgress + "', version = version + 1 where " + ID + " = ? and status = '" + Assigned + "' and " + AssignedTo + " = ?"
	DoneStatusUpdateQuery       = "update requests set status = '" + Done + "' where " + ID + " = ? and status = '" + InProgress + "' and " + AssignedTo + " = ?"
)

const (
	DashboardPath = "/dashboard"
	LoginPath     = "/login"
	GetLoginPath  = "/get-login"
	LogoutPath    = "/logout"
	HomePath      = "/"
	CreatePath    = "/create"
	StatusPath    = "/status"
	ActionPath    = "/action"
)

const (
	EmptyDBMsg       = "The database is empty"
	InternalErrorMsg = "Internal server error"
	BadRequestMsg    = "Bad request"
	BadInputMsg      = "Bad input"
	NotAllowedMsg    = "Method not allowed"
)

const (
	ID          = "id"
	Role        = "role"
	Name        = "name"
	ClientName  = "client_name"
	Phone       = "phone"
	Address     = "address"
	ProblemText = "problem_text"
	New         = "new"
	AssignedTo  = "assigned_to"
	Version     = "version"
	CreatedAt   = "created_at"
	UpdatedAt   = "updated_at"
	Master      = "master"
	Masters     = "masters"
	Status      = "status"
	User        = "user"
	Requests    = "requests"
	Action      = "action"
	Assign      = "assign"
	Assigned    = "assigned"
	Dispatcher  = "dispatcher"
	Cancel      = "cancel"
	Canceled    = "canceled"
	Start       = "start"
	InProgress  = "in_progress"
	Finish      = "finish"
	Done        = "done"
)

const (
	LoginHTML     = "login.html"
	DashboardHTML = "dashboard.html"
	CreateHTML    = "create.html"
)
