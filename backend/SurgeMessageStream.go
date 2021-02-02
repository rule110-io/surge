package surge

import (
	"encoding/binary"
	"errors"
	"io"

	"log"

	"github.com/rule110-io/surge/backend/sessionmanager"
)

// SessionWrite writes to session
func SessionWrite(Session *sessionmanager.Session, Data []byte, ID byte) (written int, err error) {

	if Session == nil || Session.Session == nil {
		return 0, errors.New("write to session error, session nil")
	}
	//Package identifier to know what we are sending
	packID := make([]byte, 1)
	packID[0] = ID

	//Create buffer of 4 bytes to put the size of the package
	buff := make([]byte, 4)
	binary.LittleEndian.PutUint32(buff, uint32(len(Data)))

	//append pack and buff
	buff = append(packID, buff...)

	//Write data
	buff = append(buff, Data...)

	//Session.session.SetWriteDeadline(time.Now().Add(60 * time.Second))
	_, err = Session.Session.Write(buff)
	if err != nil {
		return 0, err
	}

	return len(buff), err
}

//SessionRead reads from session
func SessionRead(Session *sessionmanager.Session) (data []byte, ID byte, err error) {

	headerBuffer := make([]byte, 5) //int32 size of header + 1 for packid

	// the header of 4 bytes + 1 for packid
	_, err = io.ReadFull(Session.Reader, headerBuffer)
	if err != nil {
		if err.Error() == "session closed" {
			log.Println(err)
			return nil, 0x0, err
		}
		log.Println(err)
		return nil, 0x0, err
	}

	//Get the packid
	packID := headerBuffer[0]
	log.Println(packID)

	//Get the size from the bytes
	sizeBytes := append(headerBuffer[:0], headerBuffer[1:]...)

	size := binary.LittleEndian.Uint32(sizeBytes)

	data = make([]byte, size)

	// read the full message, or return an error
	_, err = io.ReadFull(Session.Reader, data[:int(size)])
	if err != nil {
		log.Println(err)
		return nil, 0x0, err
	}

	return data, packID, err
}
