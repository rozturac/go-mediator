package mediator

//Mediator interface representations
type Mediator interface {
	Register
	Publisher
	Sender
	Behavior
}

//Mediator concrete struct
type mediator struct {
	Register
	Publisher
	Sender
	Behavior
}

//Create a new Mediator
func Create() Mediator {
	m := &mediator{}
	m.Register = newRegister()
	m.Behavior = newBehavior()
	m.Publisher = newPublisher(m.Register)
	m.Sender = newSender(m.Register, m.Behavior)
	return m
}
