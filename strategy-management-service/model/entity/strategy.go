package entity

import (
	"reflect"
	"time"
)

//ACTIONS
// Wait next step, Make a bet
type Action string

const (
	ActionWaitNextStep Action = "WAIT_NEXT_STEP"
	ActionWaitTime     Action = "WAIT_TIME"
	ActionBet          Action = "BET"
)

//OPERATION
// >, <, !=, <=, >=, ==, NONE
type operation string

const (
	OpGreater     operation = ">"
	OpLess        operation = "<"
	OpNotEq       operation = "!="
	OpLessOrEq    operation = "<="
	OpGreaterOrEq operation = ">="
	OpEquals      operation = "=="
	OpNone        operation = "NONE"
)

//CONDITIONS
// MODIFIER (param @OPERATION@ val )
// AND/OR
// param @OPERATION@ val
// ...
type union_operation string

const (
	UnOpAnd  union_operation = "AND"
	UnOpOr   union_operation = "OR"
	UnOpNone union_operation = "NONE"
)

//MODIFIER
// NOT
type modifier string

const (
	ModNot = "NOT"
)

//REPEATER
// ForEachStep
// ForEachNStep
type repeater string

const (
	RepeaterFES  = "ForEachStep"
	RepeaterFENS = "ForEachNStep"
)

// PARAMS
// isOnAdditionalPurchase
// currentWinnerId
// myId
// currentStepNumber
// stepsTillZero
// currentPrice
// currentDiscount
// timeSinceLastStep
// timeSinceLastMyBet
// stepSize
// timeSinceStart
// timeTillEnd
// minimalPrice
// acceptablePrice
// preferablePrice
// participantsCount

type paramName string

const (
	ParamIsOnAdditionalPurchase = "isOnAdditionalPurchase"
	ParamCurrentWinnerId        = "currentWinnerId"
	ParamMyId                   = "myId"
	ParamCurrentStepNumber      = "currentStepNumber"
	ParamMyCurrentBetNumber     = "myCurrentBetNumber"
	ParamStepsTillZero          = "stepsTillZero"
	ParamCurrentPrice           = "currentPrice"
	ParamCurrentDiscount        = "currentDiscount"
	ParamTimeSinceLastStep      = "timeSinceLastStep"
	ParamTimeSinceLastMyBet     = "timeSinceLastMyBet"
	ParamStepSize               = "stepSize"
	ParamTimeSinceStart         = "timeSinceStart"
	ParamTimeTillEnd            = "timeTIllEnd"
	ParamParticipantsCount      = "participantsCount"
	ParamMinimalPrice           = "minimalPrice"
	ParamAcceptablePrice        = "acceptablePrice"
	ParamPreferablePrice        = "preferablePrice"
)

type CurrentSessionState struct {
	// PARAMS
	// isOnAdditionalPurchase
	// currentWinnerId
	// userId
	// currentStepNumber
	// stepsTillZero
	// currentPrice
	// currentDiscount
	// timeSinceLastStep
	// timeSinceLastMyBet
	// stepSize
	// timeSinceStart
	// timeTillEnd
	// minimalPrice
	// acceptablePrice
	// preferablePrice
	// participantsCount
	IsOnAdditionalPurchase bool
	CurrentWinnerId        int64
	UserId                 int64
	MyCurrentBetNumber     int64
	CurrentStepNumber      int64
	StepsTillZero          int64
	CurrentPrice           float64
	CurrentDiscount        float64
	TimeSinceLastStep      time.Duration
	TimeSinceLastMyBet     time.Duration
	StepSize               float64
	TimeSinceStart         time.Duration
	TimeTillEnd            time.Duration
	ParticipantsCount      int64

	MinimalPrice    float64
	AcceptablePrice float64
	PreferablePrice float64
}

type Expression struct {
	IsParam bool
	PName   paramName

	IsConst bool
	Value   interface{}

	IsEmpty bool

	//IsCalculableExpression
}

func (e Expression) Calculate(state CurrentSessionState) interface{} {
	if !e.IsEmpty {
		if e.IsConst {
			return e.Value
		} else if e.IsParam {
			return getValueByParamName(e.PName, state)
		}
	}
	return nil
}

type Condition struct {
	M     modifier
	Param paramName
	Op    operation
	Val   Expression
}

func (r Condition) Evaluate(state CurrentSessionState) bool {
	v1 := getValueByParamName(r.Param, state)
	switch v1.(type) {
	case bool:
		return v1.(bool)
	}
	v2 := r.Val.Calculate(state)
	if reflect.TypeOf(v1) == reflect.TypeOf(v2) {
		switch v1.(type) {
		case int64:
			return r.PerformOperationInt64(v1.(int64), v2.(int64))
		case float64:
			return r.PerformOperationFloat64(v1.(float64), v2.(float64))
		case time.Duration:
			return r.PerformOperationDuration(v1.(time.Duration), v2.(time.Duration))
		}
	}
	return false
}

func (r Condition) PerformOperationInt64(v1 int64, v2 int64) bool {
	switch r.Op {
	case OpEquals:
		return v1 == v2
	case OpLess:
		return v1 < v2
	case OpLessOrEq:
		return v1 <= v2
	case OpGreater:
		return v1 > v2
	case OpGreaterOrEq:
		return v1 >= v2
	case OpNotEq:
		return v1 != v2
	}
	return false
}
func (r Condition) PerformOperationFloat64(v1 float64, v2 float64) bool {
	switch r.Op {
	case OpEquals:
		return v1 == v2
	case OpLess:
		return v1 < v2
	case OpLessOrEq:
		return v1 <= v2
	case OpGreater:
		return v1 > v2
	case OpGreaterOrEq:
		return v1 >= v2
	case OpNotEq:
		return v1 != v2
	}
	return false

}
func (r Condition) PerformOperationDuration(v1 time.Duration, v2 time.Duration) bool {
	switch r.Op {
	case OpEquals:
		return v1 == v2
	case OpLess:
		return v1 < v2
	case OpLessOrEq:
		return v1 <= v2
	case OpGreater:
		return v1 > v2
	case OpGreaterOrEq:
		return v1 >= v2
	case OpNotEq:
		return v1 != v2
	}
	return false
}

func getValueByParamName(val paramName, state CurrentSessionState) interface{} {
	switch val {
	case ParamAcceptablePrice:
		return state.AcceptablePrice
	case ParamCurrentDiscount:
		return state.CurrentDiscount
	case ParamIsOnAdditionalPurchase:
		return state.IsOnAdditionalPurchase
	case ParamMinimalPrice:
		return state.MinimalPrice
	case ParamParticipantsCount:
		return state.ParticipantsCount
	case ParamTimeSinceLastMyBet:
		return state.TimeSinceLastMyBet
	case ParamCurrentStepNumber:
		return state.CurrentStepNumber
	case ParamCurrentWinnerId:
		return state.CurrentWinnerId
	case ParamMyId:
		return state.UserId
	case ParamPreferablePrice:
		return state.PreferablePrice
	case ParamStepsTillZero:
		return state.StepsTillZero
	case ParamTimeTillEnd:
		return state.TimeTillEnd
	case ParamCurrentPrice:
		return state.CurrentPrice
	case ParamTimeSinceLastStep:
		return state.TimeSinceLastStep
	case ParamStepSize:
		return state.StepSize
	case ParamTimeSinceStart:
		return state.TimeSinceStart
	case ParamMyCurrentBetNumber:
		return state.MyCurrentBetNumber

	default:
		return nil
	}
}

//(condition AND condition) AND (condition OR condition)
//     L |operator| R             	L  |operator| R
//			L			|operator|         R
type Operator struct {
	O     union_operation
	Left  interface{}
	Right interface{}
}

func (o Operator) Evaluate(state CurrentSessionState) bool {
	switch o.O {
	case UnOpNone:
		switch o.Left.(type) {
		case Operator:
			return o.Left.(Operator).Evaluate(state)
		case Condition:
			return o.Left.(Condition).Evaluate(state)
		}
	case UnOpOr:
		leftV := false
		rightV := false
		switch o.Left.(type) {
		case Operator:
			leftV = o.Left.(Operator).Evaluate(state)
		case Condition:
			leftV = o.Left.(Condition).Evaluate(state)
		}
		switch o.Right.(type) {
		case Operator:
			rightV = o.Right.(Operator).Evaluate(state)
		case Condition:
			rightV = o.Right.(Condition).Evaluate(state)
		}
		return leftV || rightV
	case UnOpAnd:
		leftV := false
		rightV := false
		switch o.Left.(type) {
		case Operator:
			leftV = o.Left.(Operator).Evaluate(state)
		case Condition:
			leftV = o.Left.(Condition).Evaluate(state)
		}
		switch o.Right.(type) {
		case Operator:
			rightV = o.Right.(Operator).Evaluate(state)
		case Condition:
			rightV = o.Right.(Condition).Evaluate(state)
		}
		return leftV && rightV
	}
	return false
}

type ConditionSet struct {
	Op     Operator
	Action Action
	Else   *ConditionSet
}

func (c ConditionSet) Define(state CurrentSessionState) Action {
	if c.Op.Evaluate(state) {
		return c.Action
	} else {
		if c.Else != nil {
			return c.Else.Define(state)
		} else {
			return ActionWaitTime
		}
	}
}

type Strategy struct {
	// - M(Param Op Val) AND M(Param Op Val) AND ... OR ...
	//		ACTION
	//  Else
	// - M(Param Op Val) OR M(Param Op Val)
	//		ACTION
	//	Else
	//	...
	BaseConditionSet *ConditionSet
	Vars             map[string]float64
	R                repeater
	N                int
}
