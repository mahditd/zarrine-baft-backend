package controllers

import "strconv"


func parseUint(value string) uint {

	id, _ := strconv.ParseUint(value, 10, 64)

	return uint(id)
}