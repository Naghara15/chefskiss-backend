package main

import (
	"fmt"

	"chefskiss-backend/config"
)

func main() {
	fmt.Println("Mencoba koneksi ke database PostgreSQL...")

	config.ConnectDatabase()
	fmt.Println("Mantap! Test koneksi database sukses.")
}
