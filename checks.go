package flexready

// Checks defines a matrix of health checks to be run.
type Checks map[string]CheckerFunc

// CheckerFunc defines a single function for checking health of remote services.
type CheckerFunc func() error
