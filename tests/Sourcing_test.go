package tests

import (
	"github.com/dominikmayer/go-cqrs/sourcing"
	"github.com/dominikmayer/go-cqrs/tests/domain"
	"github.com/dominikmayer/go-cqrs/tests/events"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	InitLogging()
	TestingT(t)
}

// The state for the test suite
type AppTestSuite struct {
}

// Setup the test suite
var _ = Suite(&AppTestSuite{})

func (s *AppTestSuite) TestStateChangesAreRepresentedByEvents(c *C) {
	// Create a new domain object
	user := domain.NewUser("pjvds")
	c.Assert(user.Username, Equals, "pjvds")

	// We created a new user, this should be
	// captured by an event.
	c.Assert(len(user.Events()), Equals, 1)

	// Change the username of the user
	user.ChangeUsername("wwwouter")
	c.Assert(user.Username, Equals, "wwwouter")

	// We changed the username, this should be
	// captured by an event.
	c.Assert(len(user.Events()), Equals, 2)
}

func (s *AppTestSuite) TestDomainObjectCanBeBuildFromHistory(c *C) {
	// The id of our event source that we will rebuild from history.
	sourceId, _ := sourcing.ParseEventSourceId("0791d279-664d-458e-bf60-567ade140832")

	// The full history for the User domain object
	history := []sourcing.Event{
		// It was first created
		events.UserCreated{
			Username: "pjvds",
		},
		// Then the username was changed
		events.UsernameChanged{
			OldUsername: "pjvds",
			NewUsername: "wwwouter",
		},
	}

	// Create a new User domain object from history
	user := domain.NewUserFromHistory(sourceId, history)

	// It should not have the initial state.
	c.Assert(user.Username, Not(Equals), "pjvds")

	// It should have the latest state.
	c.Assert(user.Username, Equals, "wwwouter")
}

func (s *AppTestSuite) BenchmarkRebuildUserFromHistory(c *C) {
	// The full history for the User domain object
	sourceId, _ := sourcing.ParseEventSourceId("0791d279-664d-458e-bf60-567ade140832")

	// The full history for the User domain object
	history := []sourcing.Event{
		// It was first created
		events.UserCreated{
			Username: "pjvds",
		},
		// Then the username was changed
		events.UsernameChanged{
			OldUsername: "pjvds",
			NewUsername: "wwwouter",
		},
	}

	for i := 0; i < c.N; i++ {
		// Create a new User domain object from history
		domain.NewUserFromHistory(sourceId, history)
	}
}
