package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	ROCK     = 0
	PAPER    = 1
	SCISSORS = 2
)

type Game struct {
	DisplayChan chan string
	RoundChan   chan int
	Round       Round
}

type Round struct {
	RoundNumber   int
	PlayerScore   int
	ComputerScore int
}

var reader = bufio.NewReader(os.Stdin)

func (g *Game) Rounds() {
	for {
		select {
		case round := <-g.RoundChan:
			g.Round.RoundNumber += round
			g.RoundChan <- 1
		case msg := <-g.DisplayChan:
			fmt.Println(msg)
			g.DisplayChan <- ""
		}
	}
}

func (g *Game) ClearScreen() {
	if strings.Contains(runtime.GOOS, "windows") {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

func (g *Game) Intro() {
	g.DisplayChan <- fmt.Sprintf(`
Rock, Paper & Scissors
----------------------
3 Rounds game
`)
	<-g.DisplayChan
}

func (g *Game) PlayRound() bool {
	rand.Seed(time.Now().UnixNano())
	playerValue := -1

	g.DisplayChan <- fmt.Sprintf(`
Round %d
------------------
`, g.Round.RoundNumber)
	<-g.DisplayChan

	fmt.Print("Please enter rock, paper or scissors ->")
	playerChoice, _ := reader.ReadString('\n')
	playerChoice = strings.TrimSpace(playerChoice)
	computerValue := rand.Intn(3)

	if playerChoice == "rock" {
		playerValue = ROCK
	} else if playerChoice == "paper" {
		playerValue = PAPER
	} else if playerChoice == "scissors" {
		playerValue = SCISSORS
	}
	fmt.Println()
	g.DisplayChan <- fmt.Sprintf("Player chose %s", strings.ToUpper(playerChoice))
	<-g.DisplayChan
	switch computerValue {
	case ROCK:
		g.DisplayChan <- fmt.Sprintf("Computer chose ROCK")
		break
	case PAPER:
		g.DisplayChan <- fmt.Sprintf("Computer chose PAPER")
		break
	case SCISSORS:
		g.DisplayChan <- fmt.Sprintf("Computer chose SCISSORS")
		break
	default:
	}
	<-g.DisplayChan
	fmt.Println()
	if playerValue == computerValue {
		g.DisplayChan <- "DRAW"
		return false
	} else if playerValue == (computerValue+1)%3 {
		g.playerWins()
	} else {
		g.computerWins()
	}
	return true
}

func (g *Game) computerWins() {
	g.Round.ComputerScore++
	g.DisplayChan <- "Computer WINS!"
	<-g.DisplayChan

}

func (g *Game) playerWins() {
	g.Round.PlayerScore++
	g.DisplayChan <- "Player WINS!"
	<-g.DisplayChan
}

func (g *Game) PrintSummary() {
	g.DisplayChan <- fmt.Sprintf(`
Final score
-----------
Player: %d/3, Computer %d/3
`, g.Round.PlayerScore, g.Round.ComputerScore)
	<-g.DisplayChan
	if g.Round.PlayerScore > g.Round.ComputerScore {
		g.DisplayChan <- fmt.Sprintf("Player wins game!")
		<-g.DisplayChan

	} else {
		g.DisplayChan <- fmt.Sprintf("Computer wins game!")
		<-g.DisplayChan

	}
}
