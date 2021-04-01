package data

import (
	"reflect"
	"testing"
)

var targetFiles []interface{}

//func TestMain(m *testing.M) {
//	filesWithDependencies()
//	os.Exit(m.Run())
//}

func filesWithDependencies() {
	var f1 = &ShineItemInfo{}
	err := Load(filePath+"/shn/ItemInfo.shn", f1)
	if err != nil {
		log.Fatal(err)
	}

	var f2 = &ShineItemInfoServer{}
	err = Load(filePath+"/shn/ItemInfoServer.shn", f2)
	if err != nil {
		log.Fatal(err)
	}

	targetFiles = append(targetFiles, f1)
	targetFiles = append(targetFiles, f2)
}

func TestLinkedFiles(t *testing.T) {
	filesWithDependencies()

	for _, file := range targetFiles {
		t.Run(reflect.TypeOf(file).String(), func(t *testing.T) {

			f, ok := file.(FileDependency)

			if !ok {
				t.Error("file type does not implement interface FileDependency")
			}

			idfs, err := f.MissingIdentifiers(filePath)

			if err != nil {
				t.Error(err)
			}

			count := 0
			if len(idfs) > 0 {
				for k1, v1 := range idfs {
					for k2, v2 := range v1 {
						for _, v3 := range v2 {
							t.Logf("file=%v, identifier=%v, value=%v \n", k1, k2, v3)
						}
						count++
					}
				}
			}

			if count > 0 {
				t.Error("missing identifiers")
			}
		})
	}
}
