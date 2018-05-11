package statemachine

import (
	"context"
)

//ContextData because short of generic-types, use empty interface to implement
type ContextData interface{}

type Handler func(data ContextData)

type Transition struct {
	FromSate  string
	ToState   string
	Event     string
	Processor EventProcessor
}

type EventProcessor interface {
	//Before hanles before transition
	Before(ctx context.Context, t *Transition, stateHolder StateHolder, ctxData ContextData) Error
	// OnEvent handles to do all the transition
	OnEvent(ctx context.Context, t *Transition, stateHolder StateHolder, ctxData ContextData) Error
	//After handles things after transition
	After(ctx context.Context, t *Transition, stateHolder StateHolder, ctxData ContextData) Error
}

func FlowTemplate(ctx context.Context, t *Transition, stateHolder StateHolder, ctxData ContextData) Error {
	var err Error
	if err = t.Processor.Before(ctx, t, stateHolder, ctxData); err != nil {
		return err
	}
	if err = t.Processor.OnEvent(ctx, t, stateHolder, ctxData); err != nil {
		return err
	}
	if err = t.Processor.After(ctx, t, stateHolder, ctxData); err != nil {
		return err
	}
	return nil
}

type StateMachine struct {
	Name        string
	Transitions []*Transition
}
type StateHolder interface {
	CurrentState() string
}
type Locker interface {
	Lock(ctx context.Context) (StateHolder, error)
}

func NewStateMachine(name string, transitions []*Transition) *StateMachine {
	return &StateMachine{
		Name:        name,
		Transitions: transitions,
	}
}

func (m *StateMachine) Trigger(ctx context.Context, event string, locker Locker, ctxData ContextData) Error {
	stateHolder, err := locker.Lock(ctx)
	if err != nil {
		return &statemachineLockerError{locker: locker, cause: err, event: event}
	}
	var currentState = stateHolder.CurrentState()
	transition := m.Find(currentState, event)
	if transition == nil {
		return &noTransitionError{statemachine: m.Name, event: event, currentState: currentState}
	}
	return FlowTemplate(ctx, transition, stateHolder, ctxData)
}

func (m *StateMachine) Find(from string, event string) *Transition {
	for _, v := range m.Transitions {
		if v.FromSate == from && v.Event == event {
			return v
		}
	}
	return nil
}
