package model

import (
	"math/rand"
	"sync"
	"time"
)

type WeightedModel struct {
	Low  int32
	High int32
	Info string
}

func (wm *WeightedModel) Accept(random int32) bool {
	if random >= wm.Low && random < wm.High {
		return true
	}
	return false
}

type WeightedRoundRobin struct {
	Random               *rand.Rand
	WeightTotal          int32
	WeightedModelListRef []WeightedModel
	Mux                  *sync.Mutex
}

func NewWeightedRoundRobin(zones []ZoneInfo, isMark bool) WeightedRoundRobin {
	var base int32 = 0
	var tempWeightTotal int32 = 0
	tempWeightedModelList := make([]WeightedModel, 0)

	for _, zone := range zones {
		v := zone.GetRouteWeight(isMark)

		tempWeightTotal += v
		model := WeightedModel{
			Low:  base,
			High: base + v,
			Info: zone.ZoneName,
		}
		tempWeightedModelList = append(tempWeightedModelList, model)
		base = model.High
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return WeightedRoundRobin{
		Random:               r,
		WeightTotal:          tempWeightTotal,
		WeightedModelListRef: tempWeightedModelList,
		Mux:                  new(sync.Mutex),
	}
}

func (wrr *WeightedRoundRobin) GetServerAsPerAlgo() string {
	if wrr.WeightTotal > 0 {
		wrr.Mux.Lock()
		defer wrr.Mux.Unlock()
		seed := wrr.Random.Int31n(wrr.WeightTotal)
		for _, weightedModel := range wrr.WeightedModelListRef {
			if weightedModel.Accept(seed) {
				return weightedModel.Info
			}
		}
	}
	return ""
}
