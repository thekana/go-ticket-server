package custom_error

const errorCodeBase = 0

const (
	UnknownError         uint64 = errorCodeBase + 1
	InvalidJSONString    uint64 = errorCodeBase + 2
	InputValidationError uint64 = errorCodeBase + 3
	Unauthorized         uint64 = errorCodeBase + 4
	InvalidAuthToken     uint64 = errorCodeBase + 5
	AuthTokenExpired     uint64 = errorCodeBase + 6
	InactiveAccount      uint64 = errorCodeBase + 7

	DuplicateUsername uint64 = errorCodeBase + 101
)
