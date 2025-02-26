package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
)

const (
	WIDTH  = 20
	HEIGHT = 10
	LIVES  = 3
)

type Point struct {
	x, y int
}

var fixedObstacles []Point

func clearScreen() {
	cmd := exec.Command("clear")
	if os.Getenv("OS") == "Windows_NT" {
		cmd = exec.Command("cmd", "/c", "cls")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createBoard(width, height int, level int) [][]rune {
	board := make([][]rune, height+2)
	for i := range board {
		board[i] = make([]rune, width+2)
		for j := range board[i] {
			if i == 0 || i == height+1 || j == 0 || j == width+1 || (level >= 3 && i%4 == 0 && j%6 == 0) {
				board[i][j] = '#'
			} else {
				board[i][j] = ' '
			}
		}
	}
	return board
}

func printBoard(board [][]rune, score, lives, level int) {
	clearScreen()
	fmt.Printf("Životi: %d | Rezultat: %d | Level: %d\n", lives, score, level)
	for _, row := range board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}

func placeFood(board [][]rune, snake []Point) Point {
	for {
		x := rand.Intn(len(board[0])-2) + 1
		y := rand.Intn(len(board)-2) + 1
		occupied := false
		for _, p := range snake {
			if p.x == x && p.y == y {
				occupied = true
				break
			}
		}
		for _, obs := range fixedObstacles {
			if obs.x == x && obs.y == y {
				occupied = true
				break
			}
		}
		if !occupied {
			board[y][x] = '*'
			return Point{x, y}
		}
	}
}

func moveSnake(snake []Point, xSpeed, ySpeed int, grow bool) []Point {
	head := snake[0]
	newHead := Point{head.x + xSpeed, head.y + ySpeed}

	newSnake := append([]Point{newHead}, snake...)
	if !grow {
		newSnake = newSnake[:len(newSnake)-1] // Ako zmija ne raste, skida poslednji deo
	}
	return newSnake
}

func checkCollision(snake []Point, board [][]rune) bool {
	head := snake[0]
	if board[head.y][head.x] == '#' {
		return true
	}
	for _, p := range snake[1:] {
		if head == p {
			return true
		}
	}
	return false
}

func gameLoop(level, lives int) (bool, int, int) {
	width, height := WIDTH, HEIGHT
	snake := []Point{{width / 2, height / 2}}
	xSpeed, ySpeed := 1, 0
	board := createBoard(width, height, level)
	food := placeFood(board, snake)
	score := 0
	speed := time.Duration(300-level*50) * time.Millisecond
	grow := false

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
		snake = moveSnake(snake, xSpeed, ySpeed, grow)
		grow = false // Resetujemo da zmija ne raste u svakom potezu

		if checkCollision(snake, board) {
			lives--
			if lives == 0 {
				printBoard(board, score, lives, level)
				fmt.Println("Game Over!")
				return false, lives, level
			}
			return true, lives, level
		}

		// Provera da li je zmija pojela hranu
		if snake[0].x == food.x && snake[0].y == food.y {
			food = placeFood(board, snake)
			grow = true // Omogućava zmiji da naraste u sledećem potezu
			score++
		}

		board = createBoard(width, height, level)
		for _, p := range snake {
			board[p.y][p.x] = 'O'
		}
		board[food.y][food.x] = '*'
		printBoard(board, score, lives, level)
		time.Sleep(speed)

		select {
		case char := <-inputChan:
			if char == 'w' && ySpeed == 0 {
				xSpeed, ySpeed = 0, -1
			} else if char == 's' && ySpeed == 0 {
				xSpeed, ySpeed = 0, 1
			} else if char == 'a' && xSpeed == 0 {
				xSpeed, ySpeed = -1, 0
			} else if char == 'd' && xSpeed == 0 {
				xSpeed, ySpeed = 1, 0
			}
		default:
		}

		if score >= 3+level*2 {
			fmt.Println("Level passed!")
			time.Sleep(3 * time.Second)
			return true, lives, level + 1
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	level, lives := 1, LIVES
	for {
		success, newLives, newLevel := gameLoop(level, lives)
		lives = newLives
		level = newLevel
		if !success {
			break
		}
		if level > 4 {
			fmt.Println("Čestitamo! Završili ste igru!")
			break
		}
	}
}
