package flexready

// Checks defines a matrix of health checks to be run.
type Checks map[string]Checker

// AddCheck will add a health check to the matrix.
func (c Checks) AddCheck(name string, check Checker) {
	c[name] = check
}

// Checker defines the interface for checking health of remote services.
type Checker interface {
	Check() error
}

// CheckerFunc defines a single function for checking health of remote services.
type CheckerFunc func() error

func (c CheckerFunc) Check() error { return c() }
