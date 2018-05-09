package statemachine

//ContextData because short of generic-types, use empty interface to implement
type ContextData interface{}

type Handler func(data ContextData)

type State string
type Event string

type Transition struct {
	From      State
	To        State
	Event     Event
	Processor EventProcessor
}

type EventProcessor interface {
	// OnExit Event handles exiting a state
	OnExit(fromState State, ctx ContextData) Error
	// OnEvent is used to handle transitions
	OnEvent(Event Event, fromState State, toState State, ctx ContextData) Error
	// OnEnter Event handles entering a state
	OnEnter(toState State, ctx ContextData) Error
}

func FlowTemplate(t *Transition, ctx ContextData) Error {
	var err Error
	if err = t.Processor.OnExit(t.From, ctx); err != nil {
		return err
	}
	if err = t.Processor.OnEvent(t.Event, t.From, t.To, ctx); err != nil {
		return err
	}
	if err = t.Processor.OnEnter(t.To, ctx); err != nil {
		return err
	}
	return nil
}

type StateMachine struct {
	Name        string
	Transitions []*Transition
}

func NewStateMachine(name string, transitions []*Transition) *StateMachine {
	return &StateMachine{
		Name:        name,
		Transitions: transitions,
	}
}

func (m *StateMachine) Trigger(current State, event Event, ctx ContextData) Error {
	transition := m.Find(current, event)
	if transition == nil {
		return &noTransitionError{event: event, currentState: current}
	}
	return FlowTemplate(transition, ctx)
}

func (m *StateMachine) Find(from State, Event Event) *Transition {
	for _, v := range m.Transitions {
		if v.From == from && v.Event == Event {
			return v
		}
	}
	return nil
}
