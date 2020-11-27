package model

import (
	"errors"
	"strconv"
	"strings"
)

type UidRange struct {
	UidMinValue int32
	UidMaxValue int32
}

func (r *UidRange) InRange(val int32) bool {
	return val >= r.UidMinValue && val <= r.UidMaxValue
}

type UidMultiRange struct {
	ListUidRange []UidRange
}

func NewUidMultiRange() UidMultiRange {
	return UidMultiRange{
		ListUidRange: make([]UidRange, 0),
	}
}

func (mr *UidMultiRange) AddRange(uidMin, uidMax int32) {
	uidRange := UidRange{
		UidMinValue: uidMin,
		UidMaxValue: uidMax,
	}
	mr.ListUidRange = append(mr.ListUidRange, uidRange)
}

func (mr *UidMultiRange) InMultiRange(val int32) bool {
	inRange := false
	for _, rg := range mr.ListUidRange {
		if rg.InRange(val) {
			inRange = true
			break
		}
	}
	return inRange
}

func (mr *UidMultiRange) parseUidRange(multiRange, splitStr string) (UidRange, error) {
	arrayMultiRange := strings.Split(multiRange, splitStr)
	uidRange := UidRange{}
	if len(arrayMultiRange) != 2 {
		return uidRange, errors.New("invalid multiRange: " + multiRange)
	}
	if data, err := strconv.ParseInt(arrayMultiRange[0], 10, 32); err != nil {
		return uidRange, err
	} else {
		uidRange.UidMinValue = int32(data)
	}
	if data, err := strconv.ParseInt(arrayMultiRange[1], 10, 32); err != nil {
		return uidRange, err
	} else {
		uidRange.UidMaxValue = int32(data)
	}
	return uidRange, nil
}

func (mr *UidMultiRange) BuildFromString(multiRange string) error {
	if len(strings.TrimSpace(multiRange)) == 0 {
		return nil
	}
	arrayMultiRange := strings.Split(multiRange, ",")
	tempRange := make([]UidRange, 0)
	for _, oneRange := range arrayMultiRange {
		if r, err := mr.parseUidRange(oneRange, "~"); err == nil {
			tempRange = append(tempRange, r)
		} else {
			return err
		}
	}
	if len(tempRange) > 0 {
		mr.ListUidRange = tempRange
	}

	return nil
}

func (mr *UidMultiRange) CvtToUidMap(groupName string) map[int32]string {
	uidMap := make(map[int32]string)

	if len(mr.ListUidRange) == 0 {
		return uidMap
	}

	for _, uidRange := range mr.ListUidRange {
		for uid := uidRange.UidMinValue; uid < uidRange.UidMaxValue+1; uid++ {
			if uid != -1 {
				uidMap[uid] = groupName
			}
		}
	}

	return uidMap
}
