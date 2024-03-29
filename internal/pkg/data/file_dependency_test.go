package data

import (
	"reflect"
	"testing"
)

const fileDependencyUnimplemented = "file type %v does not implement interface FileDependency"

var targetFiles []interface{}

func filesWithDependencies() {
	f1 := &ShineItemInfo{}
	err := Load(filesPath+"/shn/ItemInfo.shn", f1)
	if err != nil {
		log.Fatal(err)
	}

	f2 := &ShineItemInfoServer{}
	err = Load(filesPath+"/shn/ItemInfoServer.shn", f2)
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
				t.Errorf(fileDependencyUnimplemented, reflect.TypeOf(file).String())
			}

			idfs, err := f.MissingIdentifiers(filesPath)
			if err != nil {
				t.Error(err)
			}

			count := 0
			if len(idfs) > 0 {
				for k1, v1 := range idfs {
					for k2, v2 := range v1 {
						for _, v3 := range v2 {
							t.Logf("targetFile=%v, identifier=%v, missingValue=%v \n", k1, k2, v3)
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
