package data

// Some files have linked data in other files
// Related data is linked using an identifier (ID, InxName, ItemID, ItemOptions, MobID, etc..)
// Every file may have a dependency in 0-N files
// The given Type should implement a method where dependencies on related files are checked against
type FileDependency interface {
	//// For the given file, return a list of missing indexes in linked files
	//// e.g: ItemInfoServer => ["El1", "El2"]
	//MissingIndexes(string) (map[string][]string, error)
	//// For the given file, return a list of missing indexes in linked files
	//MissingIDs(string) (map[string][]uint16, error)

	// e.g:
	// ItemInfoServer.shn => { "missingIDs": [1234,1235,1236]  }
	MissingIdentifiers(string) (Files, error)
}

type Files map[string]Identifiers

type Identifiers map[string][]interface{}
