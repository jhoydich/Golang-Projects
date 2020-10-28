package main

import (
	"github.com/google/uuid"
	"fmt"
)


func main() {
	newid := uuid.New()
	fmt.Println(newid)
}