package shn

// LoadToEDB loads shn file data to a embeddable database
func LoadToEDB(shinePath string, edbPath string) {
 // for each file in shinePath
	// switch shnFileName
	//		load file, unpack into type, persist each row to edb

	//db.Update(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte("MapInfo.shn"))
	//	err := b.Put([]byte(0), structs.Unpack(row))
	//	return err
	//})

}

// DecryptSHN decrypt binary data
func DecryptSHN(data []byte, length int) {
	if length < 1 {
		return
	}
	l := byte(length)
	for i := length - 1; i >= 0; i-- {
		var nl byte
		data[i] = data[i] ^ l
		nl = byte(i)
		nl = nl & byte(15)
		nl = nl + byte(85)
		nl = nl ^ (byte(i) * byte(11))
		nl = nl ^ l
		nl = nl ^ byte(170)
		l = nl
	}
}