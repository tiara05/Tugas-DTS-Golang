package main

import (
	"fmt"
	"os"
	"strconv"
)

// Struct untuk menyimpan data teman
type Teman struct {
	Absen     int
	Nama      string
	Alamat    string
	Pekerjaan string
	Alasan    string
}

// Database sementara teman-teman
var temanDatabase = map[int]Teman{
	1: {1, "Tiara", "Jakarta", "Developer", "Ingin mempelajari Golang Lebih Jauh"},
	2: {2, "Rahmania", "Bandung", "Designer", "Tertarik dengan kemampuan Golang"},
	3: {3, "Hadiningrum", "Surabaya", "Engineer", "Meningkatkan skill pemrograman"},
}

// Function untuk mendapatkan data teman berdasarkan absen
func getTemanByAbsen(absen int) (Teman, error) {
	teman, exists := temanDatabase[absen]
	if !exists {
		return Teman{}, fmt.Errorf("Teman dengan absen %d tidak ditemukan", absen)
	}
	return teman, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run biodata.go <absen>")
		os.Exit(1)
	}

	absen, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Absen harus berupa angka")
		os.Exit(1)
	}

	teman, err := getTemanByAbsen(absen)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Menampilkan data teman
	fmt.Println("Data Teman:")
	fmt.Printf("Absen: %d\n", teman.Absen)
	fmt.Printf("Nama: %s\n", teman.Nama)
	fmt.Printf("Alamat: %s\n", teman.Alamat)
	fmt.Printf("Pekerjaan: %s\n", teman.Pekerjaan)
	fmt.Printf("Alasan memilih kelas Golang: %s\n", teman.Alasan)
}
