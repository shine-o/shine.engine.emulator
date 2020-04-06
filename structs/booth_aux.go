package structs

//struct CHARBRIEFINFO_CAMP
//{
//  unsigned __int16 minihouse;
//  char dummy[10];
//};
type CharBriefInfoCamp struct {
	MiniHouse uint16
	Dummy     [10]byte
}

//struct STREETBOOTH_SIGNBOARD
//{
//  char signboard[30];
//};
type StreetBoothSignBoard struct {
	Text [30]byte
}
