package eventstore

import (
	"fmt"
	"github.com/pjvds/go-cqrs/sourcing"
	"github.com/pjvds/go-cqrs/tests/domain"
	"github.com/pjvds/go-cqrs/tests/events"
	. "launchpad.net/gocheck"
	"reflect"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	InitLogging()
	TestingT(t)
}

// The state for the test suite
type EventStoreTestSuite struct {
	store *EventStore
}

// Setup the test suite
var _ = Suite(&EventStoreTestSuite{})

func (s *EventStoreTestSuite) SetUpSuite(c *C) {
	register := sourcing.NewEventTypeRegister()
	namer := sourcing.NewTypeEventNamer()

	userCreatedType := reflect.TypeOf(events.UserCreated{})
	userCreatedName := namer.GetEventNameFromType(userCreatedType)
	register.Register(userCreatedName, userCreatedType)

	usernameChangedType := reflect.TypeOf(events.UsernameChanged{})
	usernameChangedName := namer.GetEventNameFromType(usernameChangedType)
	register.Register(usernameChangedName, usernameChangedType)

	store, _ := DailEventStore("http://localhost:2113", register)
	s.store = store
}

func (s *EventStoreTestSuite) TestSmoke(c *C) {
	// Create a new domain object
	user := domain.NewUser("pjvds")
	for i := 0; i < 99; i++ {
		user.ChangeUsername(fmt.Sprintf("pjvds%v", i))
	}

	state := sourcing.GetState(user)
	err := s.store.NewStream(state)
	c.Assert(err, IsNil)

	events, err := s.store.OpenStream(state.Id())
	c.Assert(err, IsNil)
	c.Assert(len(events), Equals, 100)

	c.Assert(events[0].EventSourceId, Equals, state.Id())
}
