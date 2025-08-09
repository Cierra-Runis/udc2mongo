package model

import (
	"encoding/xml"
	"strconv"
)

// See: https://www.unicode.org/reports/tr42/#d1e2752
type UCDBool bool

// See: https://www.unicode.org/reports/tr42/#lp:d1e2740
func (b *UCDBool) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case "Y":
		*b = true
	case "N":
		*b = false
	default:
		return strconv.ErrSyntax
	}
	return nil
}
