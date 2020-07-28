package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

var locker sync.Mutex

// UpdateData updates the local sector json
// func UpdateData(sectorData sector.Sector) error {
// 	// Race condition safety - only one can run at a time
// 	locker.Lock()
// 	defer locker.Unlock()

// 	sectorArray, err := ReadSectorData()
// 	if err != nil {
// 		log.Print(err)
// 	}

// 	sectorArray = append(sectorArray, sectorData)

// 	// for _, thisSector := range sectorData {
// 	// 	sectorArray = append(sectorArray, thisSector)
// 	// }

// 	data, err := json.Marshal(sectorArray)
// 	if nil != err {
// 		return err
// 	}
// 	err = ioutil.WriteFile("./tmp/data.json", data, 0700)
// 	if nil != err {
// 		return err
// 	}

// 	return nil
// }

// // ReadSectorData reads a local JSON file and returns sector.Sector array
// func ReadSectorData() ([]sector.Sector, error) {
// 	var sectorArray []sector.Sector
// 	data, err := ioutil.ReadFile("./tmp/data.json")
// 	if err != nil {
// 		return []sector.Sector{}, err
// 	}

// 	err = json.Unmarshal(data, &sectorArray)
// 	if err != nil {
// 		return []sector.Sector{}, err
// 	}
// 	return sectorArray, err
// }

// Hash a string with salt
func Hash(source string, salt string) string {
	hash := sha256.Sum256([]byte(source + salt))

	return hex.EncodeToString(hash[:])
}
