package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
)

func main() {
	// Початкові координати
	x, y := 0.5, 0.5
	// Крок переміщення
	step := 0.05

	initial := fmt.Sprintf("white\nfigure %.2f %.2f\nupdate\n", x, y)
	_, _ = http.Post("http://localhost:17000", "text/plain", bytes.NewBufferString(initial))

	time.Sleep(1 * time.Second)

	for {
		// Рух по діагоналі
		x += step
		y += step

		if x > 0.9 || y > 0.9 {
			break
		}

		script := fmt.Sprintf("move %.2f %.2f\nupdate\n", x, y)
		resp, err := http.Post("http://localhost:17000", "text/plain", bytes.NewBufferString(script))
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			resp.Body.Close()
		}

		time.Sleep(1 * time.Second)
	}
}
