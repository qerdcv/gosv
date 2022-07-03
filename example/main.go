package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/qerdcv/gosv"
)

type TestStruct struct {
	A string  `csv:"a_name"`
	B int     `csv:"b_name"`
	C int32   `csv:"c_name"`
	D int64   `csv:"d_name"`
	E float32 `csv:"e_name"`
	F float64 `csv:"f_name"`
	J bool    `csv:"j_name"`
}

func main() {
	a := make([]TestStruct, 0, 10)
	for i := 0; i < 10; i++ {
		a = append(a, TestStruct{
			A: "some value",
			B: 10 + i,
			C: -10 + int32(i),
			D: 12 + int64(i),
			E: 13.21 + float32(i),
			F: -121.31 + float64(i),
			J: i%2 == 0,
		})
	}

	buf := bytes.Buffer{}

	w := gosv.NewWriter(&buf).SetWriteHeading(true).SetDelimiter('|')
	for _, val := range a {
		w.Write(val)
	}

	b, _ := ioutil.ReadAll(&buf)
	fmt.Println(string(b))
	// Output:
	// a_name|b_name|c_name|d_name|e_name|f_name|j_name
	// some value|10|-10|12|13.21|-121.31|true
	// some value|11|-9|13|14.21|-120.31|false
	// some value|12|-8|14|15.21|-119.31|true
	// some value|13|-7|15|16.21|-118.31|false
	// some value|14|-6|16|17.21|-117.31|true
	// some value|15|-5|17|18.21|-116.31|false
	// some value|16|-4|18|19.21|-115.31|true
	// some value|17|-3|19|20.21|-114.31|false
	// some value|18|-2|20|21.21|-113.31|true
	// some value|19|-1|21|22.21|-112.31|false
}
