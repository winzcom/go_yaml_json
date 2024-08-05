package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("./yamls/health.yaml")

	bfr := bufio.NewReader(file)

	builder := BuildJSON(bfr)

	to_json, _ := json.Marshal(builder)

	fmt.Println("what is builder ", string(to_json))

}
