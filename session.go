package main

import (
	"errors"
	"github.com/emicklei/goproperties"
	"labix.org/v2/mgo"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SessionManager struct {
	configMap  map[string]properties.Properties
	sessions   map[string]*mgo.Session
	accessLock *sync.RWMutex
}

func NewSessionManager(props properties.Properties) *SessionManager {
	sess := &SessionManager{
		configMap:  make(map[string]properties.Properties),
		sessions:   make(map[string]*mgo.Session),
		accessLock: &sync.RWMutex{},
	}
	sess.SetConfig(props)
	return sess
}

func (s *SessionManager) GetAliases() []string {
	aliases := []string{}
	for k, _ := range s.configMap {
		aliases = append(aliases, k)
	}
	return aliases
}

func (s *SessionManager) Get(alias string) (*mgo.Session, bool, error) {
	config, err := s.GetConfig(alias)
	if err != nil {
		return nil, false, err
	}

	hostport := config["host"] + ":" + config["port"]
	s.accessLock.RLock()
	existing := s.sessions[hostport]
	s.accessLock.RUnlock()
	if existing != nil {
		return existing.Clone(), true, nil
	}
	s.accessLock.Lock()
	timeout := 0
	timeoutConfig := strings.Trim(config["timeout"], " ")
	if len(timeoutConfig) != 0 {
		timeout, err = strconv.Atoi(timeoutConfig)
		if err != nil {
			return nil, false, err
		}
	}
	info("connecting to [%s=%s] with timeout [%d seconds]", config["alias"], hostport, timeout)
	dialInfo := mgo.DialInfo{
		Addrs:    []string{hostport},
		Direct:   true,
		Database: config["database"],
		Username: config["username"],
		Password: config["password"],
		Timeout:  time.Duration(timeout) * time.Second,
	}
	newSession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		info("unable to connect to [%s] because:%v", hostport, err)
		newSession = nil
	} else {
		s.sessions[hostport] = newSession
	}
	s.accessLock.Unlock()
	return newSession, false, err
}

func (s *SessionManager) Close(hostport string) {
	s.accessLock.Lock()
	existing := s.sessions[hostport]
	if existing != nil {
		existing.Close()
		delete(s.sessions, hostport)
	}
	s.accessLock.Unlock()
}

func (s *SessionManager) CloseAll() {
	info("closing all sessions: ", len(s.sessions))
	s.accessLock.Lock()
	for _, each := range s.sessions {
		each.Close()
	}
	s.accessLock.Unlock()
}

func (s *SessionManager) SetConfig(props properties.Properties) {
	aliases := props.SelectProperties("mongod.*")
	for k, v := range aliases {
		parts := strings.Split(k, ".")
		alias := parts[1]
		config := s.configMap[alias]
		if config == nil {
			config = properties.Properties{}
			config["alias"] = alias
			s.configMap[alias] = config
		}
		config[parts[2]] = v
	}
}

func (s *SessionManager) GetConfig(alias string) (properties.Properties, error) {
	config := s.configMap[alias]
	if config == nil {
		return nil, errors.New("Unknown alias:" + alias)
	} else {
		return config, nil
	}
}
