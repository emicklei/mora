package session

import (
	"errors"
	"github.com/emicklei/goproperties"
	"labix.org/v2/mgo"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MongoDB Session Manager
type SessionManager struct {
	configMap  map[string]properties.Properties
	sessions   map[string]*mgo.Session
	accessLock *sync.RWMutex
}

// Creates a new Session Manager using `props` as configuration.
// For more info about properties check `mongod.*` section in `mora.properties`
func NewSessionManager(props properties.Properties) *SessionManager {
	sess := &SessionManager{
		configMap:  make(map[string]properties.Properties),
		sessions:   make(map[string]*mgo.Session),
		accessLock: &sync.RWMutex{},
	}
	sess.SetConfig(props)
	return sess
}

// Returns slice containing all configured aliases
func (s *SessionManager) GetAliases() []string {
	aliases := []string{}
	for k := range s.configMap {
		aliases = append(aliases, k)
	}
	return aliases
}

// Gets session for alias
func (s *SessionManager) Get(alias string) (*mgo.Session, bool, error) {
	// Get alias configurations
	config, err := s.GetConfig(alias)
	if err != nil {
		return nil, false, err
	}

	var uri string
	var hostport string
	var sessionId string
	if uriConfig := strings.Trim(config["uri"], " "); len(uriConfig) != 0 {
		uri = config["uri"]
		sessionId = uri
	} else {
		hostport = config["host"] + ":" + config["port"]
		sessionId = hostport
	}

	// Check if session already exists
	s.accessLock.RLock()
	existing := s.sessions[sessionId]
	s.accessLock.RUnlock()

	// Clone and return if sessions exists
	if existing != nil {
		return existing.Copy(), true, nil
	}

	// Get timeout from configuration
	s.accessLock.Lock()
	timeout := 0
	if timeoutConfig := strings.Trim(config["timeout"], " "); len(timeoutConfig) != 0 {
		timeout, err = strconv.Atoi(timeoutConfig)
		if err != nil {
			return nil, false, err
		}
	}

	// Connect to database server
	info("connecting to [%s=%s] with timeout [%d seconds]", config["alias"], sessionId, timeout)
	var newSession *mgo.Session
	if uri != "" {
		newSession, err = mgo.DialWithTimeout(uri, time.Duration(timeout)*time.Second)
	} else {
		dialInfo := mgo.DialInfo{
			Addrs:    []string{hostport},
			Direct:   true,
			Database: config["database"],
			Username: config["username"],
			Password: config["password"],
			Timeout:  time.Duration(timeout) * time.Second,
		}
		newSession, err = mgo.DialWithInfo(&dialInfo)
	}
	if err != nil {
		info("unable to connect to [%s] because:%v", sessionId, err)
		newSession = nil
	} else {
		s.sessions[sessionId] = newSession
	}
	s.accessLock.Unlock()
	return newSession, false, err
}

// Closes session based on `uri` or `host:port`
func (s *SessionManager) Close(sessionId string) {
	s.accessLock.Lock()
	if existing := s.sessions[sessionId]; existing != nil {
		existing.Close()
		delete(s.sessions, sessionId)
	}
	s.accessLock.Unlock()
}

// Closes all sessions.
func (s *SessionManager) CloseAll() {
	info("closing all sessions: ", len(s.sessions))
	s.accessLock.Lock()
	for _, each := range s.sessions {
		each.Close()
	}
	s.accessLock.Unlock()
}

// Set's session manager configuration.
func (s *SessionManager) SetConfig(props properties.Properties) {
	for k, v := range props {
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

// Get's session configurations by alias.
func (s *SessionManager) GetConfig(alias string) (properties.Properties, error) {
	if config := s.configMap[alias]; config != nil {
		return config, nil
	}
	return nil, errors.New("Unknown alias: " + alias)
}

// Log wrapper
func info(template string, values ...interface{}) {
	log.Printf("[mora] "+template+"\n", values...)
}
