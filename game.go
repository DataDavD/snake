package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
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
	oldState   *term.State
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
		input, err := reader.ReadByte()
		if err != nil {
			log.Printf("Failed to read input: %v", err)
		}
		if input == '\x1b' {
			// Potential escape sequence
			seq := make([]byte, 2)
			_, err := reader.Read(seq)
			if err == nil {
				inputQueue <- string([]byte{input, seq[0], seq[1]})
			}
		} else {
			inputQueue <- string(input)
		}
	}
}

func processInput() {
	select {
	case input := <-inputQueue:
		switch input {
		// Up
		case "w", "k", "\x1b[A":
			direction = Point{0, -1}
		// Down
		case "s", "j", "\x1b[B":
			direction = Point{0, 1}
		// Left
		case "a", "h", "\x1b[D":
			direction = Point{-1, 0}
		// Right
		case "d", "l", "\x1b[C":
			direction = Point{1, 0}
		}
	default:
	}
}

func main() {
	// Initialize random seed
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	// Start snake in the center
	snake = []Point{{width / 2, height / 2}}
	direction = Point{0, 1}
	food = Point{r.Intn(width), r.Intn(height)}

	// Set terminal to raw mode
	var err error
	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalf("Failed to set terminal to raw mode: %v", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Handle SIGINT to restore terminal mode
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(0)
	}()

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
