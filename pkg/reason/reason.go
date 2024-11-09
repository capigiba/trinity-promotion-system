package reason

import "trinity/pkg/localization"

// Define the reason constants using the LocalizedString type
var (
	// Error messages
	InvalidRequestFormat localization.LocalizedString = "error.invalid_request_format"
	InvalidRequest       localization.LocalizedString = "error.invalid_request"
	InternalServerError  localization.LocalizedString = "error.internal_server_error"
	InvalidToken         localization.LocalizedString = "error.invalid_token"
)
