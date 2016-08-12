package edgeos

import (
	"encoding/json"
	"sync"
)

// Rec holds stats on the current job
type Rec struct {
	*sync.RWMutex
	Dupes int32 `json:"dupes"`
	New   int32 `json:"new"`
	Total int32 `json:"total"`
	Uniq  int32 `json:"uniq"`
}

// Msg is a struct for relaying stats from ProcessContent
type Msg struct {
	Name string
	*Rec
}

type messenger interface {
	incDupe()
	incNew()
	getTotal()
	incUniq()
}

// GetTotal returns total records
func (m *Msg) GetTotal() {
	m.Lock()
	m.Total += m.New
	m.Unlock()
}

// IncDupe increments Dupe by 1
func (m *Msg) IncDupe() {
	m.Lock()
	m.Dupes++
	m.Unlock()
}

// IncNew increments New by 1
func (m *Msg) IncNew() {
	m.Lock()
	m.New++
	m.Unlock()
}

// IncUniq increments Uniq by 1
func (m *Msg) IncUniq() {
	m.Lock()
	m.Uniq++
	m.Unlock()
}

// NewMsg initializes a new Msg struct
func NewMsg(s string) *Msg {
	return &Msg{
		Name: s,
		Rec: &Rec{
			RWMutex: &sync.RWMutex{},
		},
	}
}

func (m *Msg) String() string {
	out, _ := json.MarshalIndent(m, "", "\t")
	return string(out)
}
