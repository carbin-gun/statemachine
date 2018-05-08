package statemachine

//ContextData because short of generic-types, use empty interface to implement
type ContextData interface{}

type Handler func(data ContextData)

type State string
type Event string

type Transition struct {
	From  State
	To    State
	Event Event
	Handler
}

type EventProcessor interface {
	// OnExit Event handles exiting a state
	OnExit(fromState State, ctx ContextData) Error
	// OnEvent is used to handle transitions
	OnEvent(Event Event, fromState State, toState State, ctx ContextData) Error
	// OnEnter Event handles entering a state
	OnEnter(toState State, ctx ContextData) Error
}

type Delegate interface {
	HandleEvent(event Event, fromState, toState State) Error
}

type DefaultDelegate struct {
	P EventProcessor
}

func (d *DefaultDelegate) HandleEvent(event Event, fromState, toState State, ctx ContextData) Error {
	if fromState != toState {
		if err := d.P.OnExit(fromState, ctx); err != nil {
			return err
		}
	}
	if err := d.P.OnEvent(event, fromState, toState, ctx); err != nil {
		return err
	}
	if fromState != toState {
		if err := d.P.OnEnter(toState, ctx); err != nil {
			return err
		}
	}
	return nil
}

type StateMachine struct {
	Name        string
	Delegate    Delegate
	Transitions []*Transition
}

func (m *StateMachine) Trigger(current State, event Event, ctx ContextData) Error {
	transition := m.Find(current, event)
	if transition == nil {
		return &noTransitionError{event: event, currentState: current}
	}
	return m.Delegate.HandleEvent(event, current, transition.To)
}

func (m *StateMachine) Find(from State, Event Event) *Transition {
	for _, v := range m.Transitions {
		if v.From == from && v.Event == Event {
			return v
		}
	}
	return nil
}
