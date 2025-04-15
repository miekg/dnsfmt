//go:generate stringer -type=tokenType

package zonefile

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

// Represents a DNS masterfile a.k.a. a zonefile
type Zonefile struct {
	entries []Entry
	suffix  []token
}

// Represents an entry in the zonefile
type Entry struct {
	tokens    []taggedToken
	IsControl bool // is this a ($INCLUDE, $TTL, ...) directoive?
	IsComment bool // is this a comment
}

// For a control entry, returns its command (e.g. $TTL, $ORIGIN, ...)
func (e Entry) Command() []byte {
	is := e.find(useControl)
	if len(is) == 0 {
		return nil
	}
	return e.tokens[is[0]].t.Value()
}

// Domain returns the ownername for the entry.
func (e Entry) Domain() []byte {
	is := e.find(useDomain)
	if len(is) == 0 {
		return nil
	}
	return e.tokens[is[0]].t.Value()
}

// Class returns the class for the entry.
func (e Entry) Class() []byte {
	is := e.find(useClass)
	if len(is) == 0 {
		return nil
	}
	return e.tokens[is[0]].t.Value()
}

// Type returns the RR Type for the entry.
func (e Entry) Type() []byte {
	is := e.find(useType)
	if len(is) == 0 {
		return nil
	}
	return e.tokens[is[0]].t.Value()
}

// TTL retrurns the TTL for the entry (if specified).
func (e Entry) TTL() *int {
	is := e.find(useTTL)
	if len(is) == 0 {
		return nil
	}
	i, _ := strconv.Atoi(string(e.tokens[is[0]].t.Value()))
	return &i
}

// Values returns the fields for the entry.
func (e Entry) Values() (ret [][]byte) {
	is := e.find(useValue)
	for i := 0; i < len(is); i++ {
		ret = append(ret, e.tokens[is[i]].t.Value())
	}
	return
}

// Comments returns the comments for the entry.
func (e Entry) Comments() (ret [][]byte) {
	is := e.find(useComment)
	for i := 0; i < len(is); i++ {
		ret = append(ret, e.tokens[is[i]].t.Value())
	}
	return
}

// Find all indices of tokens with the given use.
func (e Entry) find(use tokenUse) (is []int) {
	for i := 0; i < len(e.tokens); i++ {
		if e.tokens[i].u == use {
			is = append(is, i)
		}
	}
	return
}

type ParsingError struct {
	msg    string
	LineNo int
	ColNo  int
}

func (e *ParsingError) Error() string { return e.msg }

// List entries in the zonefile
func (z *Zonefile) Entries() (r []Entry) {
	return z.entries
}

// Write the zonefile to a bytearray
func (z *Zonefile) Save() []byte {
	var buf bytes.Buffer

	for _, e := range z.entries {
		for _, t := range e.tokens {
			buf.Write(t.t.val)
		}
	}
	for _, t := range z.suffix {
		buf.Write(t.val)
	}

	return buf.Bytes()
}

// Create a new empty zonefile
func New() (z *Zonefile) {
	return &Zonefile{}
}

// Parse bytestring containing a zonefile
func Load(data []byte) (r *Zonefile, e *ParsingError) {
	r = &Zonefile{}
	l := lex(data)

	// lex the zonefile and group tokens by line
	var line []token
	itemsInLine := 0
	for {
		t := <-l.tokens

		if t.typ == tokenEOF {
			break
		}
		if t.typ == tokenError {
			e = newParsingError(string(t.val), t)
			return
		}
		if t.IsItem() {
			itemsInLine += 1
		}

		if t.typ == tokenNewline && len(line) == 0 {
			// or empty line
			continue
		}

		line = append(line, t)
		if t.typ == tokenNewline && itemsInLine == 0 {
			// comment
			entry, err := parseLine(line)
			if err != nil {
				return nil, err
			}
			r.entries = append(r.entries, entry)
			line = nil
			itemsInLine = 0
			continue
		}

		if t.typ == tokenNewline && itemsInLine > 0 {
			entry, err := parseLine(line)
			if err != nil {
				return nil, err
			}
			r.entries = append(r.entries, entry)
			line = nil
			itemsInLine = 0
		}
	}
	if itemsInLine > 0 {
		entry, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		r.entries = append(r.entries, entry)
	} else {
		r.suffix = line
	}
	return
}

// The interesting tokens in each line are tagged by their kind so
// they are easy to find and move around.
type taggedToken struct {
	t token
	u tokenUse
}

type tokenUse int

const (
	useOther tokenUse = iota
	useType
	useClass
	useTTL
	useDomain
	useComment
	useValue
	useControl
)

// tagged token template newline
var tttNewline taggedToken = taggedToken{
	token{val: []byte{'\n'}, typ: tokenNewline}, useOther}

// tagged token template space
var tttSpace taggedToken = taggedToken{
	token{val: []byte{' '}, typ: tokenWhiteSpace}, useOther}

// tagged token template domain
var tttDomain taggedToken = taggedToken{
	token{val: []byte{'.'}, typ: tokenItem}, useDomain}

// tagged token template class
var tttClass taggedToken = taggedToken{
	token{val: []byte{'.'}, typ: tokenItem}, useClass}

// tagged token template TTL
var tttTTL taggedToken = taggedToken{
	token{val: []byte{'.'}, typ: tokenItem}, useTTL}

// tagged token template value
var tttValue taggedToken = taggedToken{
	token{val: []byte{'.'}, typ: tokenItem}, useValue}

func newParsingError(msg string, where token) *ParsingError {
	return &ParsingError{msg: msg, LineNo: where.lineno, ColNo: where.colno}
}

// Parses a tokenized line from the zonefile
func parseLine(line []token) (e Entry, err *ParsingError) {
	// add "other" tag to each token
	for _, t := range line {
		var use tokenUse
		if t.typ == tokenComment {
			use = useComment
		}
		e.tokens = append(e.tokens, taggedToken{t, use})
	}
	// Now, we figure out which item is what.  First we need to find the first item.
	iFirstItem := -1
	for i, tt := range e.tokens {
		if tt.t.IsItem() {
			iFirstItem = i
			break
		}
	}
	if iFirstItem == -1 {
		e.IsComment = e.tokens[0].t.typ == tokenComment
		return
	}

	// The first item might be a control statement, we handle that now
	if bytes.Equal(e.tokens[iFirstItem].t.Value(), []byte("$INCLUDE")) ||
		bytes.Equal(e.tokens[iFirstItem].t.Value(), []byte("$ORIGIN")) ||
		bytes.Equal(e.tokens[iFirstItem].t.Value(), []byte("$GENERATE")) ||
		bytes.Equal(e.tokens[iFirstItem].t.Value(), []byte("$TTL")) {
		e.tokens[iFirstItem].u = useControl
		e.IsControl = true
		for i := iFirstItem + 1; i < len(e.tokens); i++ {
			if e.tokens[i].t.IsItem() {
				e.tokens[i].u = useValue
			}
		}
		return
	}

	iFirstNonDomainItem := -1

	// Is there whitespace before the first item on its line?  If not,
	// then the first item is the domain and otherwise there is no domain.
	if iFirstItem == 0 || e.tokens[iFirstItem-1].t.typ == tokenNewline {
		e.tokens[iFirstItem].u = useDomain

		for i := iFirstItem + 1; i < len(e.tokens); i++ {
			if e.tokens[i].t.IsItem() {
				iFirstNonDomainItem = i
				break
			}
		}

		if iFirstNonDomainItem == -1 {
			err = newParsingError("missing type", e.tokens[iFirstItem].t)
			return
		}
	} else {
		iFirstNonDomainItem = iFirstItem
	}

	// Now, find the type item and check for the class and TTL item in between
	foundTTL, foundClass := false, false
	iType := -1
	for i := iFirstNonDomainItem; i < len(e.tokens); i++ {
		if !e.tokens[i].t.IsItem() {
			continue
		}

		// Is it a type?
		if dns_types_lut[string(e.tokens[i].t.Value())] {
			iType = i
			e.tokens[i].u = useType
			break
		}

		// A class, maybe?
		if dns_classes_lut[string(e.tokens[i].t.Value())] {
			if foundClass {
				err = newParsingError("two classes specified", e.tokens[i].t)
				return
			}
			foundClass = true
			e.tokens[i].u = useClass
			continue
		}

		// Ok, it must be a TTL
		_, err2 := strconv.Atoi(string(e.tokens[i].t.Value()))
		if err2 != nil {
			err = newParsingError(fmt.Sprintf("invalid type/class/ttl: %q", e.tokens[i].t.Value()), e.tokens[i].t)
			return
		}
		if foundTTL {
			err = newParsingError("double TTL", e.tokens[i].t)
			return
		}
		foundTTL = true
		e.tokens[i].u = useTTL
	}
	if iType == -1 {
		err = newParsingError("missing type", e.tokens[iFirstItem].t)
		return
	}

	// The remaining items are values
	for i := iType + 1; i < len(e.tokens); i++ {
		if e.tokens[i].t.IsItem() {
			e.tokens[i].u = useValue
		}
	}

	return
}

func (t token) IsItem() bool {
	return t.typ == tokenItem || t.typ == tokenQuotedItem
}

// Converts the raw data of a token to the bytestring it represents
// XXX rfc1035 isn't clear about whether e.g. "\a" makes sense;
//
//	whether "\." is interpreted allowed in quoted strings; etc
func (t token) Value() []byte {
	var what []byte
	switch t.typ {
	case tokenQuotedItem:
		what = t.val[1 : len(t.val)-1]
	case tokenItem:
		what = t.val
	default:
		return t.val
	}
	ibuf := bytes.NewBuffer(what)
	var obuf bytes.Buffer
	precedingSlash := false
	for {
		c, e := ibuf.ReadByte()
		if e != nil {
			break
		}
		if c == '\\' && !precedingSlash {
			precedingSlash = true
			continue
		}
		if precedingSlash && '0' <= c && c <= '9' {
			c2, e2 := ibuf.ReadByte()
			c3, e3 := ibuf.ReadByte()
			if e2 != nil || e3 != nil || '0' > c2 || '0' > c3 ||
				'9' < c2 || '9' < c3 {
				panic("malformed value")
			}
			v, _ := strconv.Atoi(string([]byte{c, c2, c3}))
			obuf.WriteByte(byte(v))
			continue
		}
		precedingSlash = false
		obuf.WriteByte(c)
	}
	return obuf.Bytes()
}

var dns_classes = []string{"IN", "HS", "CH"}
var dns_classes_lut map[string]bool
var dns_types = []string{}
var dns_types_lut map[string]bool

func init() {
	for k := range dns.StringToType {
		dns_types = append(dns_types, k)
	}
	dns_classes_lut = make(map[string]bool)
	for _, t := range dns_classes {
		dns_classes_lut[t] = true
	}
	// Use miekg/dns to get all record types.
	dns_types_lut = make(map[string]bool)
	for _, t := range dns_types {
		dns_types_lut[t] = true
	}
}

// Lexer
type tokenType int

const eof = 0

const (
	// Meta (zero length) tokens
	tokenError tokenType = iota
	tokenEOF

	// Non-data tokens
	tokenWhiteSpace
	tokenLeftParen
	tokenRightParen
	tokenComment

	// Data
	tokenItem
	tokenQuotedItem
	tokenNewline
)

type token struct {
	typ           tokenType // type of token
	val           []byte
	lineno, colno int // line and column number in originally parsed file
}

func (t token) String() string {
	if t.typ == tokenEOF {
		return "EOF"
	}
	return fmt.Sprintf("<%v '%v'>", t.typ, string(t.val))
}

type lexerState func(*lexer) lexerState

type lexer struct {
	buf           []byte
	pos           int
	start         int
	state         lexerState
	inGroup       bool
	tokens        chan token
	lineno        int
	colno         int
	prevLineWidth int
}

func (l *lexer) run() {
	for l.state = lexInitial; l.state != nil; {
		l.state = l.state(l)
	}
	if l.pos < len(l.buf) {
		l.errorf("could not tokenize whole file")
	}
	l.emit(tokenEOF)
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	var val []byte
	if t != tokenEOF {
		val = l.buf[l.start:l.pos]
	}
	l.tokens <- token{typ: t, val: val,
		lineno: l.lineno, colno: l.colno}
	l.start = l.pos
}

func (l *lexer) errorf(format string, args ...interface{}) lexerState {
	l.tokens <- token{typ: tokenError,
		val:    []byte(fmt.Sprintf(format, args...)),
		lineno: l.lineno, colno: l.colno}
	return nil
}

func lex(buf []byte) *lexer {
	l := &lexer{
		buf:    buf,
		tokens: make(chan token),
	}
	go l.run()
	return l
}

func (l *lexer) next() (r byte) {
	if l.pos == len(l.buf) {
		r = eof
	} else {
		r = l.buf[l.pos]
	}
	if r == '\n' {
		l.lineno += 1
		l.prevLineWidth = l.colno
		l.colno = 0
	}
	l.colno += 1
	l.pos += 1
	return
}

// backs up the lexer one byte; backup up two bytes is not allowed
func (l *lexer) backup() {
	l.pos -= 1
	l.colno -= 1
	if l.colno == 0 {
		l.lineno -= 1
		l.colno = l.prevLineWidth
	}
}

func (l *lexer) peek() byte {
	r := l.next()
	l.backup()
	return r
}

// Consumes next byte if it's in the given string
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, rune(l.next())) {
		return true
	}
	l.backup()
	return false
}

// Consumes run of bytes from the given string
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, rune(l.next())) {
	}
	l.backup()
}

// Consumes until one of the given characters if found
func (l *lexer) acceptUntil(valid string) {
	for !strings.ContainsRune(valid, rune(l.next())) {
	}
	l.backup()
}

// Start of line or after comment/item/whitespace
func lexInitial(l *lexer) lexerState {
	switch c := l.next(); {
	case c == eof:
		return nil
	case c == ' ' || c == '\t' || (l.inGroup && (c == '\n' || c == '\r')):
		return lexSpace
	case !l.inGroup && (c == '\n' || c == '\r'):
		l.emit(tokenNewline)
		return lexInitial
	case c == '"':
		return lexQuotedItem
	case c == ';':
		return lexComment
	case c == '(':
		if l.inGroup {
			return l.errorf("double (")
		}
		l.emit(tokenLeftParen)
		l.inGroup = true
		return lexInitial
	case c == ')':
		if !l.inGroup {
			return l.errorf("unexpected )")
		}
		l.emit(tokenLeftParen)
		l.inGroup = false
		return lexInitial
	default:
		return lexItem
	}
}

func lexSpace(l *lexer) lexerState {
	if l.inGroup {
		l.acceptRun(" \t\n\r")
	} else {
		l.acceptRun(" \t")
	}
	l.emit(tokenWhiteSpace)
	return lexInitial
}

func lexComment(l *lexer) lexerState {
	l.acceptUntil("\r\n\000") // XXX + eof instead of \000
	l.emit(tokenComment)
	return lexInitial
}

func lexItem(l *lexer) lexerState {
	l.acceptUntil("\r\n\t ;\000") // XXX + eof instead of \000
	l.emit(tokenItem)
	return lexInitial
}

func lexQuotedItem(l *lexer) lexerState {
	precedingSlash := false
	for {
		switch c := l.next(); {
		case c == '"' && !precedingSlash:
			l.emit(tokenQuotedItem)
			return lexInitial
		case c == '\\':
			precedingSlash = !precedingSlash
		default:
			precedingSlash = false
		}
	}
}
