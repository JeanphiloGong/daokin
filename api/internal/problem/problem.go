package problem

import "errors"

// creates the error value used by the wallet and user validation functions later in the snapshot
var Invalid = errors.New("invalid input")

// creates a shared error value representing an authentication-related problem
var Unauthenticated = errors.New("login required")

var Forbidden = errors.New("operation not allowed")

// creates a shared error value representing an operation that is not allowed
var Missing = errors.New("resource unavailable")

// creates a shared error value for an unavailable resource
var Conflict = errors.New("state or operation conflict")

// creates a shared error value representing rejection by a business rule
var Rule = errors.New("business rule rejected")

// creates a shared error value for temporarily unavailable facts
var Unavailable = errors.New("facts temporarily unavailable")
