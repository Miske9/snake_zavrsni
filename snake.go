package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

// Konstantne vrijednosti za veličinu ploče
const (
	WIDTH  = 20
	HEIGHT = 10
)

// Struktura koja predstavlja točku (x, y) na ploči
type Point struct {
	x, y int
}

// clearScreen briše ekran
func clearScreen() {
	cmd := exec.Command("clear")
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// createBoard kreira dvodimenzionalnu ploču s okvirima (#)
func createBoard(width, height int) [][]rune {
	board := make([][]rune, height+2)
	for i := range board {
		board[i] = make([]rune, width+2)
		for j := range board[i] {
			if i == 0 || i == height+1 || j == 0 || j == width+1 {
				board[i][j] = '#' // Zidovi
			} else {
				board[i][j] = ' ' // Prazan prostor
			}
		}
	}
	return board
}

// printBoard iscrtava ploču i rezultat na konzoli
func printBoard(board [][]rune, score int) {
	clearScreen()
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
	fmt.Printf("Rezultat: %d\n", score)
}

// placeFood postavlja hranu na slučajnu poziciju na ploči koja nije zauzeta zmijom ili fiksnim preprekama
func placeFood(board [][]rune, snake []Point, level int) Point {
	for {
		x := rand.Intn(len(board[0])-2) + 1
		y := rand.Intn(len(board)-2) + 1
		isOnSnake := false
		// Provjera da li se hrana postavlja na zmiju
		for _, p := range snake {
			if p.x == x && p.y == y {
				isOnSnake = true
				break
			}
		}
		if isOnSnake {
			continue
		}

		isOnObstacle := false
		// Provjera da li se hrana postavlja na fiksnu prepreku
		for _, obstacle := range fixedObstacles {
			if obstacle.x == x && obstacle.y == y {
				isOnObstacle = true
				break
			}
		}
		if isOnObstacle {
			continue
		}

		// Postavljanje hrane na praznu poziciju
		board[y][x] = '*'
		return Point{x, y}
	}
}

// moveSnake pomiče zmiju na novu poziciju prema trenutnom smjeru
func moveSnake(snake []Point, xSpeed, ySpeed int) []Point {
	head := snake[0]
	newHead := Point{head.x + xSpeed, head.y + ySpeed}
	return append([]Point{newHead}, snake[:len(snake)-1]...)
}

// checkCollision provjerava sudare zmije sa zidovima ili sobom
func checkCollision(snake []Point, board [][]rune) bool {
	head := snake[0]
	if board[head.y][head.x] == '#' { // Sudar sa zidom
		return true
	}
	for _, p := range snake[1:] {
		if head == p { // Sudar sa vlastitim tijelom
			return true
		}
	}
	return false
}

// Globalna varijabla za fiksne prepreke
var fixedObstacles []Point

// createFixedObstacles postavlja fiksne prepreke na ploču
func createFixedObstacles(board [][]rune, level int) {
	fixedObstacles = []Point{} // Resetira prepreke prije postavljanja

	numObstacles := 0
	if level == 3 {
		numObstacles = 8
	} else if level == 4 {
		numObstacles = 16
	}

	for i := 0; i < numObstacles; i++ {
		for {
			x := rand.Intn(len(board[0])-2) + 1
			y := rand.Intn(len(board)-2) + 1
			if board[y][x] == ' ' {
				board[y][x] = '#' // Prepreka
				fixedObstacles = append(fixedObstacles, Point{x, y})
				break
			}
		}
	}
}

// addFixedObstacles dodaje fiksne prepreke na ploču
func addFixedObstacles(board [][]rune) {
	for _, obstacle := range fixedObstacles {
		board[obstacle.y][obstacle.x] = '#'
	}
}

// gameLoop upravlja logikom jedne igre
func gameLoop(level int) bool {
	var width, height int
	var xSpeed, ySpeed, score int
	var snake []Point
	var food Point
	var board [][]rune
	var speed time.Duration

	switch level {
	case 1:
		width, height = 20, 10
		snake = []Point{{width / 2, height / 2}}
		xSpeed, ySpeed = 1, 0
		board = createBoard(width, height)
		food = placeFood(board, snake, level)
		speed = 300 * time.Millisecond
		score = 0
	case 2:
		width, height = 20, 10
		snake = []Point{{width / 2, height / 2}}
		xSpeed, ySpeed = 2, 0
		board = createBoard(width, height)
		food = placeFood(board, snake, level)
		speed = 200 * time.Millisecond
		score = 0
	case 3:
		width, height = 24, 12
		snake = []Point{{width / 2, height / 2}}
		xSpeed, ySpeed = 2, 0
		board = createBoard(width, height)
		food = placeFood(board, snake, level)
		if len(fixedObstacles) == 0 { // Postavlja prepreke samo jednom
			createFixedObstacles(board, level)
		}
		speed = 150 * time.Millisecond
		score = 0
	case 4:
		width, height = 28, 14
		snake = []Point{{width / 2, height / 2}}
		xSpeed, ySpeed = 3, 0
		board = createBoard(width, height)
		food = placeFood(board, snake, level)
		if len(fixedObstacles) == 0 { // Postavlja prepreke samo jednom
			createFixedObstacles(board, level)
		}
		speed = 100 * time.Millisecond
		score = 0
	}
	time.Sleep(speed)
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	inputChan := make(chan rune)
	go func() {
		for {
			if char, _, err := keyboard.GetKey(); err == nil {
				inputChan <- char
			}
		}
	}()

	for {
		snake = moveSnake(snake, xSpeed, ySpeed)

		if checkCollision(snake, board) {
			printBoard(board, score)
			fmt.Println("Game Over!")
			fmt.Println("Press 'r' to restart or 'q' to exit.")
			for char := range inputChan {
				if char == 'r' {
					fixedObstacles = []Point{} // Reset prepreka za ponovni start
					return true
				} else if char == 'q' {
					return false
				}
			}
		}

		// U gameLoop-u, uklonjen foodType i bodovanje je sada uvijek +1
		if snake[0].x == food.x && snake[0].y == food.y {
			snake = append(snake, snake[len(snake)-1]) // Dodaje dio zmije
			food = placeFood(board, snake, level)
			score++ // Uvijek dodaje 1 bod
		}

		board = createBoard(width, height)
		for _, p := range snake {
			board[p.y][p.x] = 'O'
		}
		board[food.y][food.x] = '*'

		// Dodaje fiksne prepreke na ploču
		addFixedObstacles(board)

		printBoard(board, score)

		time.Sleep(200 * time.Millisecond)

		select {
		case char := <-inputChan:
			switch char {
			case 'w':
				if ySpeed == 0 {
					xSpeed, ySpeed = 0, -1
				}
			case 's':
				if ySpeed == 0 {
					xSpeed, ySpeed = 0, 1
				}
			case 'a':
				if xSpeed == 0 {
					xSpeed, ySpeed = -1, 0
				}
			case 'd':
				if xSpeed == 0 {
					xSpeed, ySpeed = 1, 0
				}
			}
		default:
		}

		if (level == 1 && score >= 3) ||
			(level == 2 && score >= 5) ||
			(level == 3 && score >= 8) ||
			(level == 4 && score >= 10) {
			clearScreen()
			fmt.Printf("Level %d riješen!\n", level)
			time.Sleep(3 * time.Second)
			fixedObstacles = []Point{} // Reset prepreka za sljedeći level
			return true
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	level := 1
	for {
		if !gameLoop(level) {
			level = 1 // Resetiranje na prvi nivo kad igrač izgubi
			fmt.Println("Game Over!")
			fmt.Println("Press 'r' to restart or 'q' to exit.")
			var input rune
			fmt.Scanf("%c", &input)
			if input == 'r' {
				continue // Ponovno pokretanje od prvog nivoa
			} else if input == 'q' {
				break // Izlazak iz igre
			}
		} else {
			level++ // Ako je nivo završen, prelazi na sljedeći nivo
			if level > 4 {
				fmt.Println("Congratulations! You've completed the game!")
				break
			}
		}
	}
}
