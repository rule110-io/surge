package sessionmanager

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/rule110-io/surge-ui/surge/constants"

	nkn "github.com/nknorg/nkn-sdk-go"
)

var client *nkn.MultiClient
var listenFunc func(*Session)
var sessionManagerLock = sync.Mutex{}

//onConnect is a function called when a new connection is made
var onConnect func(session *Session)

//onDisconnect is a function called when a connection is lost
var onDisconnect func(addr string)

// Session is a wrapper for everything needed to maintain a surge session
type Session struct {
	Session          net.Conn
	Reader           *bufio.Reader
	lastActivityUnix int64
}

//A map to hold nkn sessions
var sessionMap map[string]*Session
var sessionLockMap map[string]*sync.Mutex

//Initialize initializes the session manager
func Initialize(nknClient *nkn.MultiClient, connectFunc func(session *Session), disconnectFunc func(addr string)) {
	sessionMap = make(map[string]*Session)
	sessionLockMap = make(map[string]*sync.Mutex)
	fileMap = make(map[string]*os.File)
	client = nknClient
	onConnect = connectFunc
	onDisconnect = disconnectFunc
}

//GetSessionLength .
func GetSessionLength() int {
	return len(sessionMap)
}

//GetSession returns a session for given address
func GetSession(Address string) (*Session, error) {
	//Check for an existing session
	lockSession(Address)
	defer unlockSession(Address)
	session, exists := sessionMap[Address]

	//create if it doesnt exist
	var err error
	if !exists {
		session, err = createSession(Address)

		if err == nil {
			sessionMap[Address] = session
		} else {
			return nil, err
		}
	}
	if exists {
		//If the sessions exists, check if its still active, if not dump it and try to create a new one.
		elapsedSinceLastActivity := time.Now().Unix() - session.lastActivityUnix
		if elapsedSinceLastActivity > 10 {
			closeSession(Address)

			session, err = createSession(Address)
			if err == nil {
				sessionMap[Address] = session
			}
		}
	}

	return session, nil
}

//GetExistingSession does not attempt to create a connection only returns existing
func GetExistingSession(Address string, timeoutInSeconds int) (*Session, bool) {
	session, exists := sessionMap[Address]

	if exists {
		//If the sessions exists, check if its still active, if not dump it and try to create a new one.
		elapsedSinceLastActivity := time.Now().Unix() - session.lastActivityUnix
		if elapsedSinceLastActivity > int64(timeoutInSeconds) {
			closeSession(Address)

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
		lastActivityUnix: time.Now().Unix(),
	}

	_, exists := sessionMap[addr]
	if exists {
		log.Println("Why are we receiving a dial when we already have a session?")
		closeSession(addr)
	}

	sessionMap[addr] = session

	go onConnect(session)

	return session
}

//UpdateActivity updates the activity timestamp on a session
func UpdateActivity(Address string) {
	session, exists := sessionMap[Address]
	if exists {
		session.lastActivityUnix = time.Now().Unix()
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
		return nil, err
	}
	reader := bufio.NewReader(nknSession)
	log.Println("Session created for: ", Address)

	session := &Session{
		Reader:           reader,
		Session:          nknSession,
		lastActivityUnix: time.Now().Unix(),
	}

	go onConnect(session)

	return session, nil
}

func closeSession(address string) {
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
	lock, exists := sessionLockMap[Addr]
	if !exists {
		lock = &sync.Mutex{}
		sessionLockMap[Addr] = lock
	}
	lock.Lock()
}

func unlockSession(Addr string) {
	lock, exists := sessionLockMap[Addr]
	if !exists {
		panic("Unlocking session lock that does not exist!")
	}
	lock.Unlock()
}
