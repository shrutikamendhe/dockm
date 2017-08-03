package dockm

// General errors.
const (
	ErrUnauthorized           = Error("Unauthorized")
	ErrResourceAccessDenied   = Error("Access denied to resource")
	ErrResourceNotFound       = Error("Unable to find resource")
	ErrUnsupportedDockerAPI   = Error("Unsupported Docker API response")
	ErrMissingSecurityContext = Error("Unable to find security details in request context")
)

// User errors.
const (
	ErrUserNotFound            = Error("User not found")
	ErrUserAlreadyExists       = Error("User already exists")
	ErrInvalidUsername         = Error("Invalid username. White spaces are not allowed.")
	ErrAdminAlreadyInitialized = Error("Admin user already initialized")
)

// Team errors.
const (
	ErrTeamNotFound      = Error("Team not found")
	ErrTeamAlreadyExists = Error("Team already exists")
)

// TeamMembership errors.
const (
	ErrTeamMembershipNotFound      = Error("Team membership not found")
	ErrTeamMembershipAlreadyExists = Error("Team membership already exists for this user and team.")
)

// ResourceControl errors.
const (
	ErrResourceControlNotFound      = Error("Resource control not found")
	ErrResourceControlAlreadyExists = Error("A resource control is already applied on this resource")
	ErrInvalidResourceControlType   = Error("Unsupported resource control type")
)

// Endpoint errors.
const (
	ErrEndpointNotFound     = Error("Endpoint not found")
	ErrEndpointAccessDenied = Error("Access denied to endpoint")
)

// Registry errors.
const (
	ErrRegistryNotFound      = Error("Registry not found")
	ErrRegistryAlreadyExists = Error("A registry is already defined for this URL")
)

// Stack errors
const (
	ErrStackNotFound = Error("Stack not found")
)

// Version errors.
const (
	ErrDBVersionNotFound = Error("DB version not found")
)

// Settings errors.
const (
	ErrSettingsNotFound = Error("Settings not found")
)

// DockerHub errors.
const (
	ErrDockerHubNotFound = Error("Dockerhub not found")
)

// Crypto errors.
const (
	ErrCryptoHashFailure = Error("Unable to hash data")
)

// JWT errors.
const (
	ErrSecretGeneration   = Error("Unable to generate secret key")
	ErrInvalidJWTToken    = Error("Invalid JWT token")
	ErrMissingContextData = Error("Unable to find JWT data in request context")
)

// File errors.
const (
	ErrUndefinedTLSFileType = Error("Undefined TLS file type")
)

// Error represents an application error.
type Error string

// Error returns the error message.
func (e Error) Error() string { return string(e) }
