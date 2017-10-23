package edgeos

import (
	"encoding/json"
	"sync"
)

// Msg is a struct for recording stats from ProcessContent
type Msg struct {
	Name string `json:"name"`
	Done bool   `json:"done"`
	*sync.RWMutex
	Dupes int `json:"dupes"`
	New   int `json:"new"`
	Total int `json:"total"`
	Uniq  int `json:"uniq"`
}

// GetTotal returns total m.New records
func (m *Msg) GetTotal() int {
	m.Lock()
	m.Total += m.New
	m.Unlock()
	return m.Total
}

// incDupe increments Dupe by 1
func (m *Msg) incDupe() {
	m.Lock()
	m.Dupes++
	m.Unlock()
}

// incNew increments New by 1
func (m *Msg) incNew() {
	m.Lock()
	m.New++
	m.Unlock()
}

// incUniq increments Uniq by 1
func (m *Msg) incUniq() {
	m.Lock()
	m.Uniq++
	m.Unlock()
}

// NewMsg initializes and returns a new Msg struct
func NewMsg(s string) *Msg {
	return &Msg{
		RWMutex: &sync.RWMutex{},
		Name:    s,
	}
}

// ReadDupes returns current value of Msg.Dupes
func (m *Msg) ReadDupes() int {
	m.RLock()
	defer m.RUnlock()
	return m.Dupes
}

// ReadNew returns current value of Msg.New
func (m *Msg) ReadNew() int {
	m.RLock()
	defer m.RUnlock()
	return m.New
}

// ReadUniq returns current value of msg.Uniq
func (m *Msg) ReadUniq() int {
	m.RLock()
	defer m.RUnlock()
	return m.Uniq
}

func (m *Msg) String() string {
	out, _ := json.MarshalIndent(m, "", "\t")
	// if err != nil {
	// 	return fmt.Sprint("cannot create string, error:", err)
	// }
	return string(out)
}
