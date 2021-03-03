package shn

import (
	"log"
	"os"
	"reflect"
	"testing"
)

var targetFiles []interface{}

func TestMain(m *testing.M) {
	filesWithDependencies()
	os.Exit(m.Run())
}

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
	for _, file := range targetFiles {
		t.Run(reflect.TypeOf(file).String(), func(t *testing.T) {
			//t.Parallel()

			f, ok := file.(FileDependency)

			if !ok {
				t.Error("file type does not implement interface FileDependency")
			}

			indexes, err := f.MissingIndexes(filePath)

			if err != nil {
				t.Error(err)
			}

			ids, err := f.MissingIDs(filePath)

			if err != nil {
				t.Error(err)
			}

			inxc := 0
			if len(indexes) > 0 {
				for k, v1 := range indexes {
					for _, v2 := range v1 {
						t.Logf("%v not in %v \n", v2, k)
						inxc++
					}
				}
			}

			idc := 0
			if len(ids) > 0 {
				for k, v1 := range ids {
					for _, v2 := range v1 {
						t.Logf("%v not in %v \n", v2, k)
						idc++
					}
				}
			}

			if idc > 0 || inxc > 0 {
				t.Errorf("missing indexes = %v, missing ids = %v", inxc, idc)
			}

		})
	}
}
