package utils

func VLDIsEmail(str string) bool {
	return RegexpContain(str, RegexpEmail)
}

func VLDIsMobile(str string) bool {
	return RegexpContain(str, RegexpMobile)
}

func VLDIsPassword(str string) bool {
	return RegexpContain(str, RegexpPassword)
}

func VLDIsIP(str string) bool {
	return RegexpContain(str, RegexpIP)
}

func VLDIsIDCard(str string) bool {
	return RegexpContain(str, RegexpIDCard)
}

func VLDIsRegexp(str string, reg string) bool {
	return RegexpContain(str, reg)
}

func VLDLenBetween(str string,min,max int) bool {
	if len(str) >= min && len(str) <= max {
		return true
	}

	return false
}





