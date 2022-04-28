package utils

import (
	"main/model/entity"
	"time"
)

var (
	AggressiveStrategy = entity.Strategy{
		BaseConditionSet: &entity.ConditionSet{
			Op: entity.Operator{
				O: entity.UnOpAnd,
				Left: entity.Condition{
					M:     "",
					Param: entity.ParamCurrentPrice,
					Op:    entity.OpGreater,
					Val: entity.Expression{
						IsParam: true,
						PName:   entity.ParamPreferablePrice,
					}, //preferablePrice
				},
				Right: entity.Condition{
					M:     "",
					Param: entity.ParamCurrentWinnerId,
					Op:    entity.OpNotEq,
					Val: entity.Expression{
						IsParam: true,
						PName:   entity.ParamMyId,
					}, // myId
				},
			},
			Action: entity.ActionBet,
			Else: &entity.ConditionSet{
				Op: entity.Operator{
					O: entity.UnOpAnd,
					Left: entity.Condition{
						M:     "",
						Param: entity.ParamCurrentPrice,
						Op:    entity.OpGreater,
						Val: entity.Expression{
							IsParam: true,
							PName:   entity.ParamAcceptablePrice,
						}, //acceptablePrice
					},
					Right: entity.Condition{
						M:     "",
						Param: entity.ParamTimeSinceLastMyBet,
						Op:    entity.OpGreater,
						Val: entity.Expression{
							IsConst: true,
							Value:   1 * time.Minute,
						}, // 1hour
					},
				},
				Action: entity.ActionBet,
				Else: &entity.ConditionSet{
					Op: entity.Operator{
						O: entity.UnOpAnd,
						Left: entity.Condition{
							M:     "",
							Param: entity.ParamCurrentPrice,
							Op:    entity.OpGreater,
							Val: entity.Expression{
								IsParam: true,
								PName:   entity.ParamMinimalPrice,
							}, //preferablePrice
						},
						Right: entity.Operator{
							O: entity.UnOpAnd,
							Left: entity.Condition{
								M:     "",
								Param: entity.ParamCurrentWinnerId,
								Op:    entity.OpNotEq,
								Val: entity.Expression{
									IsParam: true,
									PName:   entity.ParamMyId,
								}, // myId

							},
							Right: entity.Condition{
								M:     "",
								Param: entity.ParamIsOnAdditionalPurchase,
								Op:    entity.OpNone,
								Val:   entity.Expression{IsEmpty: true},
							},
						},
					},
					Action: entity.ActionBet,
					Else:   nil,
				},
			},
		},
		Vars: nil,
		//R:    "",
		//Strategy granularity (seconds)
		N: 5,
	}

	WaitingStrategy = entity.Strategy{
		BaseConditionSet: &entity.ConditionSet{
			Op: entity.Operator{
				O: entity.UnOpNone,
				Left: entity.Condition{
					M:     "",
					Param: entity.ParamMyCurrentBetNumber,
					Op:    entity.OpEquals,
					Val: entity.Expression{
						IsConst: true,
						Value:   0,
					},
				},
				Right: nil,
			},
			Action: entity.ActionBet,
			Else: &entity.ConditionSet{
				Op: entity.Operator{
					O: entity.UnOpAnd,
					Left: entity.Condition{
						M:     "",
						Param: entity.ParamCurrentPrice,
						Op:    entity.OpGreater,
						Val: entity.Expression{
							IsParam: true,
							PName:   entity.ParamMinimalPrice,
						}, // minimalPrice
					},
					Right: entity.Condition{
						M:     "",
						Param: entity.ParamIsOnAdditionalPurchase,
						Op:    entity.OpEquals,
						Val: entity.Expression{
							IsEmpty: true,
						}, //
					},
				},
				Action: entity.ActionBet,
				Else:   nil,
			},
		},
		Vars: nil,
		//R:                "",
		N: 5,
	}

	ProgressiveStrategy = entity.Strategy{
		BaseConditionSet: &entity.ConditionSet{
			Op: entity.Operator{
				O: entity.UnOpAnd,
				Left: entity.Condition{
					M:     "",
					Param: entity.ParamCurrentPrice,
					Op:    entity.OpGreater,
					Val: entity.Expression{
						IsParam: true,
						PName:   entity.ParamPreferablePrice,
					},
				},
				Right: entity.Condition{
					M:     "",
					Param: entity.ParamTimeSinceLastMyBet,
					Op:    entity.OpGreater,
					Val: entity.Expression{
						IsConst: true,
						Value:   600 * time.Second,
					},
				},
			},
			Action: entity.ActionBet,
			Else: &entity.ConditionSet{
				Op: entity.Operator{
					O: entity.UnOpAnd,
					Left: entity.Operator{
						O: entity.UnOpAnd,
						Left: entity.Condition{
							M:     "",
							Param: entity.ParamCurrentPrice,
							Op:    entity.OpGreater,
							Val: entity.Expression{
								IsParam: true,
								PName:   entity.ParamAcceptablePrice,
							},
						},
						Right: entity.Condition{
							M:     "",
							Param: entity.ParamCurrentPrice,
							Op:    entity.OpLess,
							Val: entity.Expression{
								IsParam: true,
								PName:   entity.ParamPreferablePrice,
							},
						},
					},
					Right: entity.Condition{
						M:     "",
						Param: entity.ParamTimeSinceLastMyBet,
						Op:    entity.OpGreater,
						Val: entity.Expression{
							IsConst: true,
							Value:   100 * time.Second,
						},
					},
				},
				Action: entity.ActionBet,
				Else: &entity.ConditionSet{
					Op: entity.Operator{
						O: entity.UnOpAnd,
						Left: entity.Operator{
							O: entity.UnOpAnd,
							Left: entity.Condition{
								M:     "",
								Param: entity.ParamCurrentPrice,
								Op:    entity.OpGreater,
								Val: entity.Expression{
									IsParam: true,
									PName:   entity.ParamMinimalPrice,
								},
							},
							Right: entity.Condition{
								M:     "",
								Param: entity.ParamCurrentPrice,
								Op:    entity.OpLess,
								Val: entity.Expression{
									IsParam: true,
									PName:   entity.ParamAcceptablePrice,
								},
							},
						},
						Right: entity.Condition{
							M:     "",
							Param: entity.ParamTimeSinceLastMyBet,
							Op:    entity.OpGreater,
							Val: entity.Expression{
								IsConst: true,
								Value:   5 * time.Second,
							},
						},
					},
					Action: entity.ActionBet,
					Else:   nil,
				},
			},
		},
		Vars: nil,
		//R:                "",
		N: 1,
	}
)
