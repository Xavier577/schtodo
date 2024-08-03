package objects

import (
	"encoding/json"
	"errors"

	"github.com/Xavier577/schtodo/pkg/typings"
)

var (
	ErrTargetCannotBeNil   = errors.New("TARGET CANNOT BE NIL")
	ErrDestCannotBeNil     = errors.New("DEST CANNOT BE NIL")
	ErrTargetMustBeAStruct = errors.New("TARGET MUST BE A STRUCT")
	ErrDestMustBeAStruct   = errors.New("DEST MUST BE A STRUCT")
)

func MarshalStructToMap(target any) (map[string]any, error) {

	var objMap map[string]any

	if target == nil {
		return objMap, nil
	}

	if !typings.IsStruct(target) && !typings.IsPointerToStruct(target) {
		return objMap, ErrTargetMustBeAStruct
	}

	var jsonObj []byte
	var errJsonMarshal error

	if typings.IsStruct(target) {
		jsonObj, errJsonMarshal = json.Marshal(&target)
	} else {

		jsonObj, errJsonMarshal = json.Marshal(target)
	}

	if errJsonMarshal != nil {
		return nil, errJsonMarshal
	}

	errJsonUnmarshal := json.Unmarshal(jsonObj, &objMap)

	if errJsonUnmarshal != nil {
		return nil, errJsonUnmarshal
	}

	return objMap, nil
}

func MarshalMapToStruct(target map[string]any, dest any) error {

	if target == nil {
		return ErrTargetCannotBeNil
	}

	if typings.IsZeroValue(dest) {
		return ErrDestCannotBeNil
	}

	if !typings.IsStruct(dest) && !typings.IsPointerToStruct(dest) {
		return ErrDestMustBeAStruct
	}

	targetJson, errJsonMarshal := json.Marshal(target)

	if errJsonMarshal != nil {
		return errJsonMarshal
	}

	var errJsonUnmarshal error

	if typings.IsStruct(dest) {
		errJsonUnmarshal = json.Unmarshal(targetJson, &dest)
	} else {
		errJsonUnmarshal = json.Unmarshal(targetJson, dest)
	}

	if errJsonUnmarshal != nil {
		return errJsonUnmarshal
	}

	return nil
}

func MustMarshalStructMerge(target any, dests ...any) {
	errStructMerge := MarshalStructMerge(target, dests...)

	if errStructMerge != nil {
		panic(errStructMerge)
	}
}

// MarshalStructMerge merges the common fields from two or more structs.
// It required the structs both the target and dest structs to have json tags
//
// Note: the zero values of the dest fields would be ignored
func MarshalStructMerge(target any, dests ...any) error {

	if !typings.IsStruct(target) && !typings.IsPointerToStruct(target) {

		return ErrTargetMustBeAStruct
	}

	var targetMap map[string]any
	var errMarshalTargetToMap error

	if typings.IsStruct(target) {
		targetMap, errMarshalTargetToMap = MarshalStructToMap(&target)
	} else {
		targetMap, errMarshalTargetToMap = MarshalStructToMap(target)

	}

	if errMarshalTargetToMap != nil {
		return errMarshalTargetToMap
	}

	for _, dest := range dests {

		if !typings.IsStruct(dest) && !typings.IsPointerToStruct(dest) && !typings.IsMap(dest) {
			return ErrDestMustBeAStruct

		}

		var destMap map[string]any
		var errMarshalDestToMap error

		if typings.IsStruct(dest) {
			destMap, errMarshalDestToMap = MarshalStructToMap(&dest)
		} else if typings.IsPointerToStruct(dest) {
			destMap, errMarshalDestToMap = MarshalStructToMap(dest)
		} else {
			destMap = dest.(map[string]any)
		}

		if errMarshalDestToMap != nil {
			return errMarshalDestToMap
		}

		for key, val := range destMap {
			if !typings.IsZeroValue(val) {
				_, ok := targetMap[key]

				if ok {
					targetMap[key] = val
				}
			}
		}

	}

	return MarshalMapToStruct(targetMap, target)
}
