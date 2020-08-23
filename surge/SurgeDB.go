package surge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/xujiajun/nutsdb"
)

const fileBucketName = "fileBucket"
const settingBucketName = "settingsBucket"

var db *nutsdb.DB

//InitializeDb initializes db
func InitializeDb() {
	var err error

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	opt := nutsdb.DefaultOptions

	if runtime.GOOS == "darwin" {
		dir, _ = os.UserHomeDir()
		dir = dir + string(os.PathSeparator) + ".surge"
		opt.Dir = dir + string(os.PathSeparator) + "db"
	} else {
		opt.Dir = dir + string(os.PathSeparator) + "db"
	}

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
			e, err := tx.Get(fileBucketName, fileKey)
			if err != nil {
				return err
			}

			err = json.Unmarshal(e.Value, result)
			return err
		}); err != nil {
		return nil, err
	}

	return result, nil
}

func dbGetAllFiles() []File {
	files := []File{}

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(fileBucketName)
			if err != nil {
				return err
			}

			for _, entry := range entries {

				newFile := &File{}
				json.Unmarshal(entry.Value, newFile)
				files = append(files, *newFile)
			}

			return nil
		}); err != nil {
		log.Println(err)
	} else {
		return files
	}
	return files
}

func dbDeleteFile(Hash string) error {
	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			key := []byte(Hash)
			if err := tx.Delete(fileBucketName, key); err != nil {
				return err
			}
			return nil
		}); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//CloseDb .
func CloseDb() {
	db.Close()
}

//DbWriteSetting .
func DbWriteSetting(Name string, value string) error {
	err := db.Update(
		func(tx *nutsdb.Tx) error {

			keyBytes := []byte(Name)
			valueBytes := []byte(value)

			if err := tx.Put(settingBucketName, keyBytes, valueBytes, 0); err != nil {
				return err
			}
			return nil
		})
	return err
}

//DbReadSetting .
func DbReadSetting(Name string) (string, error) {
	result := ""
	key := []byte(Name)

	if err := db.View(
		func(tx *nutsdb.Tx) error {
			bytes, err := tx.Get(settingBucketName, key)
			if err != nil {
				return err
			}

			result = string(bytes.Value)

			return err
		}); err != nil {
		return "", err
	}

	return result, nil
}
