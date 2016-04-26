package game

import (
	"fmt"
	"strings"
)

// Move is an action the player can take
type Move int

const (
	Down Move = iota
	Left
	Right
	TurnLeft
	TurnRight
	Skip
)

func serializeMoves(mvs *[]Move) (str string, err error) {
	mvToStr := map[Move]string{
		Down:      "down",
		Left:      "left",
		Right:     "right",
		TurnLeft:  "turnleft",
		TurnRight: "turnright",
		Skip:      "skip",
	}
	strs := []string{}
	for _, mv := range *mvs {
		strs = append(strs, mvToStr[mv])
	}
	return strings.Join(strs, ","), nil
}

// Player can send instructions to the game and get instructions back
type Player struct {
	input  <-chan string
	output chan<- string
}

// Process starts a goroutine processing the player. It returns two channels.
// the first you can read from to get the current state. Whenever you recieve
// a state, you should send a list of moves on the second channel and then
// wait for a new state again.
func (p *Player) Process() (<-chan *State, chan<- *[]Move) {
	sts := make(chan *State)
	mvss := make(chan *[]Move)
	go func() {
		st := NewState()
		for {
			fmt.Printf("Waiting on mvs %v\n", mvss)
			// wait for a message from the server or some moves to send to it
			// from the user
			select {
			case msg := <-p.input:
				fmt.Printf("Got a message; %v\n", msg)
				shouldMove, err := st.processLine(msg)
				if err != nil {
					panic(err)
				}
				if shouldMove {
					sts <- st
				}
			case mvs := <-mvss:
				fmt.Print("Got a Move\n")

				msg, err := serializeMoves(mvs)
				if err != nil {
					panic(err)
				}
				p.output <- msg
			}
		}
	}()
	fmt.Printf("Sending mvs %v\n", mvss)
	return sts, mvss
}
