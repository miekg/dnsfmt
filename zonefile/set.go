package zonefile

import (
	"bytes"
	"errors"
)

func (t *token) SetValue(v []byte) {
	if !t.IsItem() {
		panic("not implemented") // XXX
	}
	if bytes.IndexByte(v, ' ') >= 0 {
		// XXX replace non-printable characters (even though the rfc
		//     would allow them).
		tmp := bytes.Replace(v, []byte("\\"), []byte("\\\\"), -1)
		tmp = bytes.Replace(v, []byte("\""), []byte("\\\""), -1)
		t.typ = tokenQuotedItem
		t.val = []byte("\"" + string(tmp) + "\"")
		return
	}
	tmp := bytes.Replace(v, []byte("\\"), []byte("\\\\"), -1)
	tmp = bytes.Replace(v, []byte("\""), []byte("\\\""), -1)
	t.typ = tokenItem
	t.val = tmp
	return
}

// Set the the ith value of the entry
func (e *Entry) SetValue(i int, v []byte) error {
	if len(v) == 0 {
		return errors.New("value must be non-empty")
	}
	is := e.find(useValue)
	if len(is) <= i {
		return errors.New("index of value is too high")
	}
	e.tokens[is[i]].t.SetValue(v)
	return nil
}

// Changes the domain in the entry
func (e *Entry) SetDomain(v []byte) error {
	if e.IsControl {
		return errors.New("control entry does not have a domain")
	}
	is := e.find(useDomain)

	if len(is) == 1 {
		// If there is a domain item, simply change its value
		if len(v) != 0 {
			e.tokens[is[0]].t.SetValue(v)
			return nil
		}

		//  ... or delete it if we don't want a domain
		e.tokens = append(e.tokens[:is[0]], e.tokens[is[0]+1:]...)
	}

	// There is no domain and we don't want one, so that's ok!
	if len(v) == 0 {
		return nil
	}

	// If there is no domain item in the entry, add it
	iFirstToken := e.startOfLine()
	tDomain := tttDomain
	tDomain.t.SetValue(v)
	toAdd := []taggedToken{tDomain}
	if e.tokens[iFirstToken].t.typ != tokenWhiteSpace {
		toAdd = append(toAdd, tttSpace)
	}
	e.tokens = append(e.tokens[:iFirstToken], append(toAdd, e.tokens[iFirstToken:]...)...)
	return nil
}

// startOfLine finds the first token on the main line of the entry.
func (e Entry) startOfLine() (r int) {
	var firstItem int
	for i := 0; i < len(e.tokens); i++ {
		if e.tokens[i].t.IsItem() {
			firstItem = i
			break
		}
	}
	for i := firstItem; i >= 0; i-- {
		if e.tokens[i].t.typ == tokenNewline {
			r = i + 1
			return
		}
	}
	return 0
}
