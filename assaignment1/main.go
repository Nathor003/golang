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
	data := make([]dataOrang, 3)

	fmt.Println("index", index)

	data[0] = dataOrang{name: "andi", alamat: "jakarta", pekerjaan: "sales", alasan: "Penasaran aja"}
	data[1] = dataOrang{name: "ando", alamat: "jakarta", pekerjaan: "programmer", alasan: "butuh kerja"}
	data[2] = dataOrang{name: "budi", alamat: "jakarta", pekerjaan: "telemarketing", alasan: "pindah haluan"}

	idx, _ := strconv.Atoi(index)

	fmt.Println(data[idx])
}
