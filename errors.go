package nacoswrap

import "errors"

// Sentinel errors returned by the client when an underlying client is unavailable.
var (
	ErrNamingNotInit = errors.New("nacoswrap: naming client not initialized")
	ErrConfigNotInit = errors.New("nacoswrap: config client not initialized")
)
