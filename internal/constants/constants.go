package constants

type requestIDCtxKey struct{}

var RequestIDKey = requestIDCtxKey{}

const UserContext = "user"
const RequestIDHeader = "X-Request-Id"
