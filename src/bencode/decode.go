package bencode

import (
	"errors"
	"strconv"
)

// Torrent struct
type Torrent struct {
	data []byte
	pos  int
}

// Any type
type Any interface{}

// NewTorrent constructor
func NewTorrent(s string) Torrent {
	var bs = []byte(s)
	var t = Torrent{bs, 0}
	return t
}

func (t *Torrent) eat(c byte) error {
	if t.data[t.pos] != c {
		return errors.New("input data error")
	}
	t.pos = t.pos + 1
	return nil
}

func (t *Torrent) next() (byte, error) {
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

func (t *Torrent) watch() (byte, error) {
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

func (t *Torrent) int() (int, error) {
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

func (t *Torrent) number() (int, error) {
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

func (t *Torrent) string() (string, error) {
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

func (t *Torrent) list() (Any, error) {
	var (
		list []Any
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

func (t *Torrent) dict() (Any, error) {
	m := make(map[string]Any)
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

func (t *Torrent) value() (Any, error) {
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
