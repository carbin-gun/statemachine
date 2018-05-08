package statemachine

import "fmt"

//Error is an error when processing event and state changing
type Error interface {
	error
	Event() string
	CurrentState() string
}

var _ Error = new(noTransitionError)

type noTransitionError struct {
	event        Event
	currentState State
}

func (e *noTransitionError) Error() string {
	return fmt.Sprintf("statemachine error: no transition for event [%s] at current state [%s]", e.event, e.currentState)
}

func (e *noTransitionError) Event() string {
	return string(e.event)
}
func (e *noTransitionError) CurrentState() string {
	return string(e.currentState)
}

type TransitionError struct {
	cause   error
	event   Event
	current State
}

func (e *TransitionError) Error() string {
	return fmt.Sprintf("statemachine transition error: current state [%s] event [%s] but error:%v happened", e.current, e.event, e.cause)
}
func (e *TransitionError) Event() string {
	return string(e.event)
}
func (e *TransitionError) CurrentState() string {
	return string(e.current)
}
