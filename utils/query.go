package utils

import "strconv"

func StringToInt(param []string) ([]int, error) {

	var len = len(param)
	var intParam = make([]int, len)

	for i := 0; i < len; i++ {
		intValue, err := strconv.Atoi(param[i])
		if err != nil {
			return nil, err
		}
		intParam[i] = intValue
	}

	return intParam, nil
}
