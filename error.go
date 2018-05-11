package statemachine

import "fmt"

type ErrorCode uint8

const (
	NoTransition    ErrorCode = iota
	LockerError
	TransitionError
)

type Error interface {
	error
	CurrentState() string
	Event() string
	Code() ErrorCode
}
type noTransitionError struct {
	statemachine string
	event        string
	currentState string
}

func (e *noTransitionError) Error() string {
	return fmt.Sprintf("statemachine transition error:  statemachine [%s]no transition for event [%s] at current state [%s]", e.statemachine, e.event, e.currentState)
}

func (e *noTransitionError) CurrentState() string {
	return e.currentState
}
func (e *noTransitionError) Event() string {
	return e.event
}
func (e *noTransitionError) Code() ErrorCode {
	return NoTransition
}

type statemachineLockerError struct {
	locker Locker
	cause  error
	event  string
}

func (e *statemachineLockerError) Error() string {
	return fmt.Sprintf("statemachine lock error,locker:%v,event:%v,error:%v", e.locker, e.event, e.cause)
}
func (e *statemachineLockerError) CurrentState() string {
	return "unknown"
}
func (e *statemachineLockerError) Event() string {
	return e.event
}
func (e *statemachineLockerError) Code() ErrorCode {
	return LockerError
}

type transitionError struct {
	cause   error
	event   string
	current string
}

func NewTransitionError(cause error, event string, currentState string) Error {
	return &transitionError{
		cause:   cause,
		event:   event,
		current: currentState,
	}
}

func (e *transitionError) Error() string {
	return fmt.Sprintf("statemachine transition error: current state [%s] event [%s] but error:%v happened", e.current, e.event, e.cause)
}

func (e *transitionError) CurrentState() string {
	return e.current
}
func (e *transitionError) Event() string {
	return e.event
}
func (e *transitionError) Code() ErrorCode {
	return TransitionError
}
