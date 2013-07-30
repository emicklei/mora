package main

import (
	"github.com/emicklei/goproperties"
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
	info("closing all sessions:", len(sessions))
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

func openSession(config properties.Properties) (*mgo.Session, bool, error) {
	hostport := config["host"] + ":" + config["port"]
	sessionAccessMutex.RLock()
	existing := sessions[hostport]
	sessionAccessMutex.RUnlock()
	if existing != nil {
		return existing.Clone(), true, nil
	}
	sessionAccessMutex.Lock()
	info("connecting to [%s=%s]", config["alias"], hostport)
	dialInfo := mgo.DialInfo{
		Addrs:    []string{hostport},
		Direct:   true,
		Database: config["database"],
		Username: config["username"],
		Password: config["password"],
	}
	newSession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		info("unable to connect to [%s] because:%v", hostport, err)
		newSession = nil
	} else {
		sessions[hostport] = newSession
	}
	sessionAccessMutex.Unlock()
	return newSession, false, err
}
