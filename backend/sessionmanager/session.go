package sessionmanager

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/rule110-io/surge/backend/constants"

	nkn "github.com/nknorg/nkn-sdk-go"
)

var client *nkn.MultiClient
var listenFunc func(*Session)
var sessionManagerLock = sync.Mutex{}

//onConnect is a function called when a new connection is made, isDialIn is whether its dial out or dialing in
var onConnect func(session *Session, isDialIn bool)

//onDisconnect is a function called when a connection is lost
var onDisconnect func(addr string)

// Session is a wrapper for everything needed to maintain a surge session
type Session struct {
	Session          net.Conn
	Reader           *bufio.Reader
	LastActivityUnix int64
}

//A map to hold nkn sessions
var sessionMap map[string]*Session
var sessionLockMap map[string]*sync.Mutex
var sessionLockMapLock sync.Mutex

//Initialize initializes the session manager
func Initialize(nknClient *nkn.MultiClient, connectFunc func(session *Session, isDialIn bool), disconnectFunc func(addr string)) {
	sessionMap = make(map[string]*Session)
	sessionLockMap = make(map[string]*sync.Mutex)
	sessionLockMapLock = sync.Mutex{}
	//fileMap = make(map[string]*os.File)
	client = nknClient
	onConnect = connectFunc
	onDisconnect = disconnectFunc
}

//GetSessionLength .
func GetSessionLength() int {
	return len(sessionMap)
}

//GetSessionsString .
func GetSessionsString() string {
	arr := []string{}
	for k := range sessionMap {
		arr = append(arr, k)
	}
	return strings.Join(arr, ",")
}

//GetSession returns a session for given address
func GetSession(Address string, timeoutInSeconds int) (*Session, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic: %+v\n", r)
		}
	}()

	//Check for an existing session

	session, exists := sessionMap[Address]

	//create if it doesnt exist
	var err error
	if !exists {
		session, err = createSession(Address)

		if err == nil {
			sessionLockMapLock.Lock()
			sessionMap[Address] = session
			sessionLockMapLock.Unlock()
		} else {
			return nil, err
		}
	}
	/*
		if exists {
			//If the sessions exists, check if its still active, if not dump it and try to create a new one.
			elapsedSinceLastActivity := time.Now().Unix() - session.LastActivityUnix
			if elapsedSinceLastActivity > int64(timeoutInSeconds) {
				closeSession(Address)

				session, err = createSession(Address)
				if err == nil {
					sessionLockMapLock.Lock()
					sessionMap[Address] = session
					sessionLockMapLock.Unlock()
				}
			}
		}*/

	return session, nil
}

//GetExistingSession does not attempt to create a connection only returns existing
func GetExistingSession(Address string, timeoutInSeconds int) (*Session, bool) {
	session, exists := sessionMap[Address]

	if exists {
		//If the sessions exists, check if its still active, if not dump it and try to create a new one.
		elapsedSinceLastActivity := time.Now().Unix() - session.LastActivityUnix
		if elapsedSinceLastActivity > int64(timeoutInSeconds) {
			closeSession(Address)

			return nil, false
		}
	}

	return session, exists
}

//GetExistingSessionWithoutClosing does not attempt to create a connection only returns existing
func GetExistingSessionWithoutClosing(Address string, timeoutInSeconds int) (*Session, bool) {
	session, exists := sessionMap[Address]

	if exists {
		//If the sessions exists, check if its still active, if not dump it and try to create a new one.
		elapsedSinceLastActivity := time.Now().Unix() - session.LastActivityUnix
		if elapsedSinceLastActivity > int64(timeoutInSeconds) {
			return nil, false
		}
	}

	return session, exists
}

//CloseSession handles session termination and removes from map
/*func CloseSession(address string) {
	lockSession(Address)
	defer unlockSession(Address)

	closeSession(address)
}*/

//AcceptSession accepts a incoming session connection
func AcceptSession(acceptedConnection net.Conn) *Session {
	addr := acceptedConnection.RemoteAddr().String()

	listenReader := bufio.NewReader(acceptedConnection)
	session := &Session{
		Reader:           listenReader,
		Session:          acceptedConnection,
		LastActivityUnix: time.Now().Unix(),
	}

	//Give it a 10 sec headstart, old session workers take up to 10 sec to timeout, then to fetch the new session this would then already be timedout.
	session.LastActivityUnix = time.Now().Unix() + constants.WorkerGetSessionTimeout
	sessionMap[addr] = session

	go onConnect(session, true)

	return session
}

//UpdateActivity updates the activity timestamp on a session
func UpdateActivity(Address string) {
	sessionLockMapLock.Lock()
	defer sessionLockMapLock.Unlock()
	session, exists := sessionMap[Address]

	if exists {
		session.LastActivityUnix = time.Now().Unix()
	}
}

func createSession(Address string) (*Session, error) {
	sessionConfig := nkn.GetDefaultSessionConfig()
	sessionConfig.MTU = 16384
	//sessionConfig.CheckTimeoutInterval = 1
	//sessionConfig.InitialRetransmissionTimeout = 1
	//sessionConfig.MaxRetransmissionTimeout = 1

	dialConfig := &nkn.DialConfig{
		SessionConfig: sessionConfig,
		DialTimeout:   constants.NknClientDialTimeout,
	}

	nknSession, err := client.DialWithConfig(Address, dialConfig)
	if err != nil {
		log.Println("Failed to create a session with ", Address, err)
		fmt.Println(string("\033[31m"), "Failed to create a session with ", Address, err, string("\033[0m"))

		//If we have a session that didnt come in after dial
		/*acceptedSession, dialupExists := GetExistingSession(Address, constants.NknClientDialTimeout)
		if dialupExists {
			fmt.Println(string("\033[31m"), "but inbound (accepted) dialup was received in the meantime", Address, err, string("\033[0m"))
			return acceptedSession, nil
		}*/

		return nil, err
	}
	reader := bufio.NewReader(nknSession)
	log.Println("Session created for: ", Address)

	session := &Session{
		Reader:           reader,
		Session:          nknSession,
		LastActivityUnix: time.Now().Unix(),
	}

	go onConnect(session, false)

	return session, nil
}

func IsExistingSession(address string) bool {
	sessionLockMapLock.Lock()
	defer sessionLockMapLock.Unlock()
	_, exists := sessionMap[address]
	return exists
}

func closeSession(address string) {
	sessionLockMapLock.Lock()
	defer sessionLockMapLock.Unlock()

	session, exists := sessionMap[address]

	//Close nkn session, nill out the pointers
	if exists {
		if session.Session != nil {
			session.Session.Close()
		}
		session.Session = nil
		session.Reader = nil
	}
	session = nil

	//Delete from the map
	delete(sessionMap, address)

	log.Println("Download Session closed for: ", address)
	fmt.Println(string("\033[31m"), "Download Session closed for: ", address, string("\033[0m"))

	go onDisconnect(address)
}

func lockSession(Addr string) {
	sessionLockMapLock.Lock()
	defer sessionLockMapLock.Unlock()
	lock, exists := sessionLockMap[Addr]
	if !exists {
		lock = &sync.Mutex{}
		sessionLockMap[Addr] = lock
	}
	lock.Lock()
}

func unlockSession(Addr string) {
	sessionLockMapLock.Lock()
	defer sessionLockMapLock.Unlock()
	lock, exists := sessionLockMap[Addr]
	if !exists {
		panic("Unlocking session lock that does not exist!")
	}
	lock.Unlock()
}
