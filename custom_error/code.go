package custom_error

const errorCodeBase = 0

const (
	UnknownError         uint64 = errorCodeBase + 1
	InvalidJSONString    uint64 = errorCodeBase + 2
	InputValidationError uint64 = errorCodeBase + 3
	Unauthorized         uint64 = errorCodeBase + 4
	InvalidAuthToken     uint64 = errorCodeBase + 5
	AuthTokenExpired     uint64 = errorCodeBase + 6
	UserNotFound         uint64 = errorCodeBase + 7
	DuplicateUsername    uint64 = errorCodeBase + 8
	ConcurrencyIssue     uint64 = errorCodeBase + 9
	InsufficientQuota    uint64 = errorCodeBase + 10
	NotEnoughPrivileges  uint64 = errorCodeBase + 11
	DBError              uint64 = errorCodeBase + 12
	RedisError           uint64 = errorCodeBase + 13
	BadInput             uint64 = errorCodeBase + 14
	EventNotFound        uint64 = errorCodeBase + 15
)
