package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	width  = 50
	height = 10
)

type Point struct {
	x, y int
}

var (
	snake      []Point
	direction  Point
	food       Point
	gameOver   bool
	inputQueue = make(chan string)
)

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to clear the screen: %v", err)
	}
}

func draw() {
	clearScreen()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if contains(snake, Point{x, y}) {
				fmt.Print("O")
			} else if food.x == x && food.y == y {
				fmt.Print("*")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func contains(slice []Point, point Point) bool {
	for _, p := range slice {
		if p == point {
			return true
		}
	}
	return false
}

func move() {
	head := snake[0]
	newHead := Point{head.x + direction.x, head.y + direction.y}

	if newHead.x < 0 || newHead.x >= width || newHead.y < 0 || newHead.y >= height || contains(snake, newHead) {
		gameOver = true
		return
	}

	snake = append([]Point{newHead}, snake...)

	if newHead == food {
		food = Point{rand.Intn(width), rand.Intn(height)}
	} else {
		snake = snake[:len(snake)-1]
	}
}

func inputListener() {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		inputQueue <- input
	}
}

func processInput() {
	select {
	case input := <-inputQueue:
		switch input {
		// Up
		case "w\n", "k\n", "\x1b[A":
			direction = Point{0, -1}
		// Down
		case "s\n", "j\n", "\x1b[B":
			direction = Point{0, 1}
		// Left
		case "a\n", "h\n", "\x1b[D":
			direction = Point{-1, 0}
		// Right
		case "d\n", "l\n", "\x1b[C":
			direction = Point{1, 0}
		}
	default:
	}
}

func main() {
	// init
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// Start snake in the center
	snake = []Point{{width / 2, height / 2}}
	direction = Point{0, 1}
	food = Point{r.Intn(width), r.Intn(height)}

	go inputListener()
	ticker := time.NewTicker(500 * time.Millisecond)

	for !gameOver {
		select {
		case <-ticker.C:
			processInput()
			move()
			draw()
		}
	}

	fmt.Println("Game Over!")
}
