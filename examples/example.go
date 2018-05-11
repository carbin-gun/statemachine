package main

import (
	"github.com/carbin-gun/statemachine"
	"log"
	"fmt"
	"context"
)

func main() {
	transitions := []*statemachine.Transition{
		{FromSate: "New", ToState: "Paid", Event: "Payment", Processor: new(OrderPayProcessor)},
		{FromSate: "Paid", ToState: "InPrepare", Event: "Confirm", Processor: new(OrderConfirmProcessor)},
	}
	var sm = statemachine.NewStateMachine("订单状态机", transitions)
	var locker = &OrderLocker{ID: 123}
	var ctx = context.Background()
	err := sm.Trigger(ctx, "Payment", locker, &OrderPayContextData{OrderNo: 123, Price: 100})
	if err != nil {
		log.Fatalf("状态机切换异常:%v", err)
	}
	err =sm.Trigger(ctx, "Confirm", locker, &OrderConfirmContextData{Operator: 123})
	if err != nil {
		log.Fatalf("状态机切换异常:%v", err)
	}
	log.Printf("订单当前状态为:%v\n",storageOrder.Status)
}

var storageOrder = &Order{ID: 123, Status: "New"}

type OrderLocker struct {
	ID int64
}
type Order struct {
	ID     int64
	Status string
}

func (o *Order) CurrentState() string {
	return o.Status
}
func (l *OrderLocker) Lock(ctx context.Context) (statemachine.StateHolder, error) {
	//这里对数据进行加锁.模拟存储在数据库中的数据
	return storageOrder, nil
}

type OrderConfirmContextData struct {
	OrderNo  int64
	Operator int64
}

type OrderConfirmProcessor struct {
}

func (p *OrderConfirmProcessor) Before(ctx context.Context, t *statemachine.Transition, stateHolder statemachine.StateHolder, ctxData statemachine.ContextData) statemachine.Error {
	fmt.Printf("exit state:%v\n", t.FromSate)
	return nil

}

// OnEvent is used to handle transitions
func (p *OrderConfirmProcessor) OnEvent(ctx context.Context, t *statemachine.Transition, stateHolder statemachine.StateHolder, ctxData statemachine.ContextData) statemachine.Error {
	fmt.Printf("transform from [%v] to [%v] by event:%s,data:%#v\n", t.FromSate, t.ToState, t.Event, ctxData)
	//模拟更改存储
	storageOrder.Status = t.ToState
	return nil
}

// OnEnter Event handles entering a state
func (p *OrderConfirmProcessor) After(ctx context.Context, t *statemachine.Transition, stateHolder statemachine.StateHolder, ctxData statemachine.ContextData) statemachine.Error {
	fmt.Printf("Enter new State :%v\n", t.ToState)
	return nil
}

type OrderPayContextData struct {
	OrderNo int64
	Price   int64
}

type OrderPayProcessor struct {
}

func (p *OrderPayProcessor) Before(ctx context.Context, t *statemachine.Transition, stateHolder statemachine.StateHolder, ctxData statemachine.ContextData) statemachine.Error {
	fmt.Printf("exit state:%v\n", t.FromSate)
	return nil
}

// OnEvent is used to handle transitions
func (p *OrderPayProcessor) OnEvent(ctx context.Context, t *statemachine.Transition, stateHolder statemachine.StateHolder, ctxData statemachine.ContextData) statemachine.Error {
	fmt.Printf("transform from [%v] to [%v] by event:%s,data:%#v\n", t.FromSate, t.ToState, t.Event, ctxData)
	//模拟更改存储
	storageOrder.Status = t.ToState
	return nil

}

// OnEnter Event handles entering a state
func (p *OrderPayProcessor) After(ctx context.Context, t *statemachine.Transition, stateHolder statemachine.StateHolder, ctxData statemachine.ContextData) statemachine.Error {
	fmt.Printf("Enter new State :%v\n", t.ToState)
	return nil
}
