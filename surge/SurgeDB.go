package surge

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xujiajun/nutsdb"
)

const fileBucketName = "fileBucket"

var db *nutsdb.DB

//InitializeDb initializes db
func InitializeDb() {
	opt := nutsdb.DefaultOptions
	opt.Dir = "./db"

	var err error
	db, err = nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	dbGetAllFiles()

}

func dbInsertFile(File File) {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {

			fileKey := []byte(File.FileHash)
			fileBytes, _ := json.Marshal(File)

			if err := tx.Put(fileBucketName, fileKey, fileBytes, 0); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Fatal(err)
	}
}

func dbGetFile(Key string) (*File, error) {
	result := &File{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			fileKey := []byte(Key)

			if e, err := tx.Get(fileBucketName, fileKey); err != nil {
				return err
			} else {
				json.Unmarshal(e.Value, result)
				return err
			}
		}); err != nil {
		return nil, err
	}
	return result, nil
}

func dbGetAllFiles() {
	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(fileBucketName)
			if err != nil {
				return err
			}

			for _, entry := range entries {

				newFile := &File{}
				json.Unmarshal(entry.Value, newFile)
				fmt.Println(string(entry.Key), newFile)
			}

			return nil
		}); err != nil {
		log.Println(err)
	}
}

//CloseDb .
func CloseDb() {
	db.Close()
}
