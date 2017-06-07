package dht

import (
	"bytes"
	"errors"
	"strconv"
)

// decodeStatus struct
type decodeStatus struct {
	data []byte
	pos  int
}

// Any type
type Any interface{}

// Map struct
type Map map[string]Any

// List stuct
type List []Any

// Decode string
func Decode(b []byte) (Any, error) {
	var t = decodeStatus{b, 0}
	return t.value()
}

func (t *decodeStatus) eat(c byte) error {
	if t.data[t.pos] != c {
		return errors.New("input data error")
	}
	t.pos = t.pos + 1
	return nil
}

func (t *decodeStatus) next() (byte, error) {
	var (
		b   byte
		err error
	)
	b, err = t.watch()
	if err != nil {
		return b, err
	}
	t.pos = t.pos + 1
	return b, nil
}

func (t *decodeStatus) watch() (byte, error) {
	var b byte
	var end = t.pos + 1
	if end > len(t.data) {
		return b, errors.New("data out of index")
	}
	var bs = t.data[t.pos : t.pos+1]
	return bs[0], nil
}

func isNumber(b byte) bool {
	if b < '0' || b > '9' {
		return false
	}
	return true
}

func (t *decodeStatus) int() (int, error) {
	var start = t.pos
	for {
		var b, err = t.watch()
		if err != nil {
			return 0, err
		}
		if !isNumber(b) {
			if t.pos == start {
				var err = errors.New("int not found")
				return 0, err
			}
			return strconv.Atoi(string(t.data[start:t.pos]))
		}
		t.pos = t.pos + 1
	}
}

func (t *decodeStatus) number() (int, error) {
	var (
		err error
		i   int
	)
	err = t.eat('i')
	if err != nil {
		return i, err
	}

	i, err = t.int()
	if err != nil {
		return i, err
	}

	err = t.eat('e')
	if err != nil {
		return i, err
	}
	return i, nil
}

func (t *decodeStatus) string() (string, error) {
	var (
		length int
		err    error
		s      string
	)
	length, err = t.int()
	if err != nil {
		return s, err
	}
	err = t.eat(':')
	if err != nil {
		return s, err
	}
	if length <= 0 {
		return s, errors.New("parse string length error")
	}
	s = string(t.data[t.pos : t.pos+length])
	t.pos = t.pos + length
	return s, nil
}

func (t *decodeStatus) list() (List, error) {
	var (
		list List
		err  error
		v    Any
		b    byte
	)
	err = t.eat('l')
	if err != nil {
		return list, err
	}
	for {
		v, err = t.value()
		if err != nil {
			return list, err
		}
		list = append(list, v)
		b, err = t.watch()
		if err != nil {
			return list, err
		}
		if b == 'e' {
			t.pos = t.pos + 1
			break
		}
	}
	return list, err
}

func (t *decodeStatus) dict() (Map, error) {
	m := make(Map)
	var (
		err error
		k   string
		v   Any
		b   byte
	)
	err = t.eat('d')
	if err != nil {
		return m, err
	}

	for {
		k, err = t.string()
		if err != nil {
			return m, err
		}

		v, err = t.value()
		if err != nil {
			return m, err
		}

		m[k] = v
		b, err = t.watch()
		if err != nil {
			return m, err
		}
		if b == 'e' {
			t.pos = t.pos + 1
			break
		}
	}
	return m, err
}

func (t *decodeStatus) value() (Any, error) {
	var (
		guessByte byte
		err       error
	)
	guessByte, err = t.watch()
	if err != nil {
		return nil, err
	}
	switch true {
	case isNumber(guessByte):
		return t.string()
	case guessByte == 'i':
		return t.number()
	case guessByte == 'l':
		return t.list()
	case guessByte == 'd':
		return t.dict()
	default:
		return nil, errors.New("parse value type error")
	}
}

type encodeStatus struct {
	buf bytes.Buffer
}

// Encode Value
func Encode(v Any) ([]byte, error) {
	var (
		e   encodeStatus
		err error
	)
	err = e.valueTo(v)
	return e.toBytes(), err
}

func (e *encodeStatus) toString() string {
	return e.buf.String()
}

func (e *encodeStatus) toBytes() []byte {
	return e.buf.Bytes()
}

func (e *encodeStatus) valueTo(v Any) error {
	switch v.(type) {
	case int:
		i, _ := v.(int)
		return e.numberTo(i)
	case string:
		s, _ := v.(string)
		return e.stringTo(s)
	case List:
		l, _ := v.(List)
		return e.listTo(l)
	case Map:
		m, _ := v.(Map)
		return e.mapTo(m)
	}
	return nil
}

func (e *encodeStatus) mapTo(m Map) error {
	if len(m) <= 0 {
		return errors.New("empty map error")
	}
	e.buf.WriteString("d")
	for k, v := range m {
		e.stringTo(k)
		e.valueTo(v)
	}
	e.buf.WriteString("e")
	return nil
}

func (e *encodeStatus) listTo(l List) error {
	if len(l) <= 0 {
		return errors.New("empty list error")
	}
	e.buf.WriteString("l")
	for _, value := range l {
		e.valueTo(value)
	}
	e.buf.WriteString("e")
	return nil
}

func (e *encodeStatus) stringTo(s string) error {
	length := len(s)
	e.buf.WriteString(strconv.Itoa(length))
	e.buf.WriteString(":")
	e.buf.WriteString(s)
	return nil
}

func (e *encodeStatus) numberTo(i int) error {
	e.buf.WriteString("i")
	s := strconv.Itoa(i)
	e.buf.WriteString(s)
	e.buf.WriteString("e")
	return nil
}
