package utils

import "strings"

func NormalizePhone(phone string) string {

	phone = strings.TrimSpace(phone)

	if strings.HasPrefix(phone, "+98") {
		phone = "0" + phone[3:]
	}

	return phone
}
