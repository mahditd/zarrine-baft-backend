package controllers

import "strconv"

func parseUint(value string) (uint, error) {

	id, err := strconv.ParseUint(value, 10, 64)

	if err != nil {
		return 0, err
	}

	return uint(id), nil
}
