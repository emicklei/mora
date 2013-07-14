package main

import (
	"labix.org/v2/mgo"
	"sync"
)

var sessions map[string]*mgo.Session
var sessionAccessMutex *sync.RWMutex

func init() {
	sessions = map[string]*mgo.Session{}
	sessionAccessMutex = new(sync.RWMutex)
}

func closeSessions() {
	sessionAccessMutex.Lock()
	for _, each := range sessions {
		each.Close()
	}
	sessionAccessMutex.Unlock()
}

func closeSession(hostport string) {
	sessionAccessMutex.Lock()
	existing := sessions[hostport]
	if existing != nil {
		existing.Close()
		delete(sessions, hostport)
	}
	sessionAccessMutex.Unlock()
}

// hostport like localhost:27017
func openSession(hostport string) (*mgo.Session, error) {
	sessionAccessMutex.RLock()
	existing := sessions[hostport]
	sessionAccessMutex.RUnlock()
	if existing != nil {
		return existing, nil
	}
	sessionAccessMutex.Lock()
	newSession, err := mgo.Dial(hostport)
	sessions[hostport] = newSession
	sessionAccessMutex.Unlock()
	return newSession, err
}
