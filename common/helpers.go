package common

import (
	"encoding/json"
	"fmt"
)

func StructToMap(myStruct any) (map[string]any, error) {
	jsonString, err := json.Marshal(myStruct)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to JSON: %s", err)
	}
	var myMap map[string]any
	err = json.Unmarshal(jsonString, &myMap)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to JSON: %s", err)
	}

	return myMap, nil
}
