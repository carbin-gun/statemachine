package main

import (
	"github.com/carbin-gun/statemachine"
	"log"
	"fmt"
)

func main() {
	transitions := []*statemachine.Transition{
		{From: "New", To: "Paid", Event: "Payment", Processor: new(OrderPayProcessor)},
		{From: "Paid", To: "InPrepare", Event: "Confirm", Processor: new(OrderConfirmProcessor)},
	}
	var sm = statemachine.NewStateMachine("订单状态机", transitions)
	err := sm.Trigger("New", "Payment", &OrderPayContextData{OrderNo: 123, Price: 100})
	if err != nil {
		log.Fatalf("状态机切换异常:%v", err)
	}
}

type OrderConfirmContextData struct {
	OrderNo  int64
	Operator int64
}

type OrderConfirmProcessor struct {
}

func (p *OrderConfirmProcessor) OnExit(fromState statemachine.State, ctx statemachine.ContextData) statemachine.Error {

	return nil

}

// OnEvent is used to handle transitions
func (p *OrderConfirmProcessor) OnEvent(Event statemachine.Event, fromState statemachine.State, toState statemachine.State, ctx statemachine.ContextData) statemachine.Error {

	return nil
}

// OnEnter Event handles entering a state
func (p *OrderConfirmProcessor) OnEnter(toState statemachine.State, ctx statemachine.ContextData) statemachine.Error {

	return nil
}

type OrderPayContextData struct {
	OrderNo int64
	Price   int64
}

type OrderPayProcessor struct {
}

func (p *OrderPayProcessor) OnExit(fromState statemachine.State, ctx statemachine.ContextData) statemachine.Error {
	fmt.Printf("exit state:%v\n", fromState)
	return nil
}

// OnEvent is used to handle transitions
func (p *OrderPayProcessor) OnEvent(Event statemachine.Event, fromState statemachine.State, toState statemachine.State, ctx statemachine.ContextData) statemachine.Error {
	fmt.Printf("transform from [%v] to [%v] by event:%s,data:%#v\n", fromState, toState, Event, ctx)
	return nil

}

// OnEnter Event handles entering a state
func (p *OrderPayProcessor) OnEnter(toState statemachine.State, ctx statemachine.ContextData) statemachine.Error {
	fmt.Printf("Enter new State :%v", toState)
	return nil
}
