package main

import "fmt"

type User struct {
	Id             int
	Name, Location string
}

func (u *User) Greetings() string {
	return fmt.Sprintf("Hi %s from %s",
		u.Name, u.Location)
}

type Player struct {
	User
	GameId int
}

func main() {
	p := &Player{
		GameId: 1,
	}
	fffffuck := p
	p.Id = 42
	p.Name = "Matt"
	p.Location = "LA"

	fffffuck.User = User{
		Id:       1,
		Name:     "fuck",
		Location: "fuckkke",
	}

	// fffffuck.Name = "ssss"
	fmt.Println(p.Greetings())
	fmt.Printf("GameId %v", fffffuck.GameId)
}
