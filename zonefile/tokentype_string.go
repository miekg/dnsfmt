// Code generated by "stringer -type=tokenType"; DO NOT EDIT

package zonefile

import "fmt"

func (i tokenType) String() string {
	index := [...]uint8{0, 10, 18, 33, 47, 62, 74, 83, 98, 110}
	const tokens = "tokenErrortokenEOFtokenWhiteSpacetokenLeftParentokenRightParentokenCommenttokenItemtokenQuotedItemtokenNewline"
	if i < 0 || i >= tokenType(len(index)-1) {
		return fmt.Sprintf("tokenType(%d)", i)
	}
	return tokens[index[i]:index[i+1]]
}
