package undo

import (
	"fmt"
	"net/url"
)

type Entry struct {
	url.Values
}

type Stack struct {
	Current int
	Entries []*Entry
}

func NewEntry(form url.Values) *Entry {
	form.Set("_cached", "true")
	return &Entry{form}
}

func (s *Stack) Replace(e *Entry) {
	switch {
	case len(s.Entries) == 0:
		s.Entries = []*Entry{e}
	case s.Current > 0:
		s.Current -= 1
		fallthrough
	default:
		s.Entries[s.Current] = e
	}
}

func (s *Stack) Push(e *Entry) {
	switch {
	case s.Current < len(s.Entries):
		s.Entries = append(s.Entries[:s.Current], e)
		s.Current += 1
	default:
		s.Entries = append(s.Entries, e)
		s.Current += 1
	}
}

func (s *Stack) Pop() *Entry {
	switch {
	case s.Current == 0:
		s.Entries = nil
		return nil
	default:
		e := s.Entries[s.Current-1]
		s.Entries = s.Entries[:s.Current]
		s.Current -= 1
		return e
	}
}

func (s *Stack) Top() *Entry {
	if s.Current < len(s.Entries) {
		return s.Entries[s.Current]
	}
	return nil
}

func (s *Stack) Back() *Stack {
	s.Current -= 1
	if s.Current < 0 {
		s.Current = 0
	}
	return s
}

func (s *Stack) Forward() *Stack {
	s.Current += 1
	if l := len(s.Entries); s.Current > l {
		s.Current = l
	}
	return s
}

func (s *Stack) String() string {
	r := fmt.Sprintf("Current: %d\n", s.Current)
	for i, e := range s.Entries {
		r += fmt.Sprintf("Entry %d: %#v\n", i, e)
	}
	return r
}

func (s *Stack) Count() int {
	return len(s.Entries)
}
