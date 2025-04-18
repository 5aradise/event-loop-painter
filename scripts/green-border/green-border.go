package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	script := `green
bgrect 0.05 0.05 0.95 0.95
update`

	resp, err := http.Post("http://localhost:17000", "text/plain", bytes.NewBufferString(script))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Response status:", resp.Status)
}
