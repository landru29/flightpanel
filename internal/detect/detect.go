package detect

import (
	"regexp"
)

var SerialPattern = regexp.MustCompile(`(?i)^\\Device\\USBSER`)
