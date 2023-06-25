package db

const (
	DBName       = "hotel_reservation"
	TestDBName   = "test_hotel_reservation"
	MongoURI     = ""
	TestMongoURI = ""
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
