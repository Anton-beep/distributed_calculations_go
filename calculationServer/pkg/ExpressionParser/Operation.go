package ExpressionParser

import (
	"errors"
	"sync"
	"time"
)

const (
	ADD = iota
	SUBTRACT
	DIVIDE
	MULTIPLY
)

type Operation struct {
	operator         int
	chan1            <-chan int
	chan2            <-chan int
	execTimeAdd      time.Duration
	execTimeSubtract time.Duration
	execTimeDivide   time.Duration
	execTimeMultiply time.Duration
	lock             sync.RWMutex
}

func NewOperation(operator int, chan1, chan2 chan int, execTimeAdd, execTimeSubtract, execTimeDivide, execTimeMultiply time.Duration) *Operation {
	return &Operation{
		operator:         operator,
		chan1:            chan1,
		chan2:            chan2,
		execTimeAdd:      execTimeAdd,
		execTimeSubtract: execTimeSubtract,
		execTimeDivide:   execTimeDivide,
		execTimeMultiply: execTimeMultiply,
	}
}

func (o *Operation) SetExecTimes(execTimeAdd, execTimeSubtract, execTimeDivide, execTimeMultiply time.Duration) error {
	if execTimeAdd < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	if execTimeSubtract < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	if execTimeDivide < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	if execTimeMultiply < 0 {
		return errors.New("execution time cannot be smaller than 0")
	}
	o.execTimeAdd = execTimeAdd
	o.execTimeSubtract = execTimeSubtract
	o.execTimeDivide = execTimeDivide
	o.execTimeMultiply = execTimeMultiply
	return nil
}

func (o *Operation) Calculate() (chan int, error) {
	num1, num2 := <-o.chan1, <-o.chan2
	res := 0
	if o.operator == ADD {
		o.lock.RLock()
		d := o.execTimeAdd
		o.lock.RUnlock()
		time.Sleep(d)

		res = num1 + num2
	} else if o.operator == SUBTRACT {
		o.lock.RLock()
		d := o.execTimeSubtract
		o.lock.RUnlock()
		time.Sleep(d)

		res = num1 - num2
	} else if o.operator == DIVIDE {
		o.lock.RLock()
		d := o.execTimeDivide
		o.lock.RUnlock()
		time.Sleep(d)

		if num2 == 0 {
			return nil, errors.New("division by zero")
		}
		res = num1 / num2
	} else if o.operator == MULTIPLY {
		o.lock.RLock()
		d := o.execTimeMultiply
		o.lock.RUnlock()
		time.Sleep(d)

		res = num1 * num2
	} else {
		return nil, errors.New("unknown/wrong Operator")
	}

	chanOut := make(chan int)
	go func() {
		chanOut <- res
	}()

	return chanOut, nil
}
