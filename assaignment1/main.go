package main

import (
	"fmt"
	"os"
	"strconv"
)

type dataOrang struct {
	name      string
	alamat    string
	pekerjaan string
	alasan    string
}

func main() {

	index := os.Args[1]

	person := []dataOrang{
		{
			name: "andi", alamat: "jakarta", pekerjaan: "sales", alasan: "Penasaran aja",
		},
		{
			name: "ando", alamat: "jakarta", pekerjaan: "programmer", alasan: "butuh kerja",
		},
		{
			name: "budi", alamat: "jakarta", pekerjaan: "telemarketing", alasan: "pindah haluan",
		},
	}

	idx, _ := strconv.Atoi(index)

	fmt.Println(person[idx])
}
