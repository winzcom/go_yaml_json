package main

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

func TestY2J(t *testing.T) {
	// test service works fine
	t.Run("success test ", func(t *testing.T) {
		file, _ := os.Open("./yamls/test.yaml")

		defer file.Close()

		btr := bufio.NewReader(file)
		builder := BuildJSON(btr)

		if reflect.TypeOf(builder).Kind() != reflect.Map {
			t.Fatal("failed to get approprate field from map")
		}

	})
}
