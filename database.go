package main

type ClientProfile struct {
	Id    string
	Name  string
	Gmail string
	Token string
}

var database = map[string]ClientProfile{
	"user1": {
		Id:    "user1",
		Name:  "John Doe",
		Gmail: "johndoe@gmail.com",
		Token: "john123",
	},
	"user2": {
		Id:    "user2",
		Name:  "Michael",
		Gmail: "michael@gmail.com",
		Token: "michael123",
	},
}
