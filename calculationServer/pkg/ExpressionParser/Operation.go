package ExpressionParser

import (
	"sync"
)

const (
	ADD = iota
	SUBTRACT
	DIVIDE
	MULTIPLY
)

type Operation struct {
	operator       int
	chan1          <-chan int
	chan2          <-chan int
	execTimeConfig ExecTimeConfig
	lock           sync.RWMutex
}

func NewOperation(operator int, chan1, chan2 chan int, execTimeConfig ExecTimeConfig) *Operation {
	return &Operation{
		operator:       operator,
		chan1:          chan1,
		chan2:          chan2,
		execTimeConfig: execTimeConfig,
		lock:           sync.RWMutex{},
	}
}

func (o *Operation) SetExecTimes(execTimeConfig ExecTimeConfig) error {
	if _, err := IsExecTimeConfigCorrect(execTimeConfig); err != nil {
		return err
	}
	o.execTimeConfig = execTimeConfig
	return nil
}

/*
func (o *Operation) CalculateRPNData() (chan int, error) {
	num1, num2 := <-o.chan1, <-o.chan2
	res := 0
	if o.operator == ADD {
		o.lock.RLock()
		d := o.execTimeConfig.ExecTimeAdd
		o.lock.RUnlock()
		time.Sleep(d)

		res = num1 + num2
	} else if o.operator == SUBTRACT {
		o.lock.RLock()
		d := o.execTimeConfig.TimeSubtract
		o.lock.RUnlock()
		time.Sleep(d)

		res = num1 - num2
	} else if o.operator == DIVIDE {
		o.lock.RLock()
		d := o.execTimeConfig.execTimeDivide
		o.lock.RUnlock()
		time.Sleep(d)

		if num2 == 0 {
			return nil, errors.New("division by zero")
		}
		res = num1 / num2
	} else if o.operator == MULTIPLY {
		o.lock.RLock()
		d := o.execTimeConfig.TimeMultiply
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
*/
