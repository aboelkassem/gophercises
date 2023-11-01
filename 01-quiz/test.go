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
	// & create pointer
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

	// * is invert of &
	// p is already a pointer, when add this sign again
	// *p = take a copy, don't use the pointer (the same object)
	notPointer := *p

	notPointer.User.Name = "not pointer"

	newPlayer := Player{
		GameId: 1,
		User: User{
			Name: "New player",
		},
	}

	newNewPlayer := newPlayer
	newNewPlayer.User.Name = "New new player"

	newPlayer.User.Name = "New player"

	// fffffuck.Name = "ssss"
	fmt.Println(p.Greetings())
	fmt.Printf("GameId: %v\n", fffffuck.GameId)
	fmt.Printf("New new player: %v", newNewPlayer.User.Name)
}
