/***
 *
 *
 * csv - https://golang-blog.blogspot.com/2020/06/csv-package-in-golang.html
 **/
package main

import(
	"fmt"
	"time"
	_ "math"
	"encoding/csv"
    "io"
    "log"
    "strings"
)

var (
	gRndSeed uint32 = 1 // последнее случайное число
	was int = 0 // была ли вычислена пара чисел
	r float64= 0 // предыдущее число
)
// Начиная с некоторого целого числа x0 =/= 0, задаваемого при помощи фукнции SRnd(),
// при каждом вызове функции Rnd() происходит вычисление нового псевдослучайного 
// числа на основе предыдущего.
func SRnd64(seed int64) {
	SRnd(uint32(seed))
}

func SRnd(seed uint32) {
	if seed == uint32(0) {
		gRndSeed = uint32(1)
	} else {
		gRndSeed = seed
	}
}
// Метод генерации случайных чисел основанный на эффекет переполнения 32-разрядных целых чисел
// возвращает равномерно распределённое случайное число
func RndU() uint32 {
	gRndSeed = gRndSeed * uint32(1664525) + uint32(1013904223)
	return gRndSeed
}
// генерировать челое число из диапазона
// с типами надо подумать...
func RndBetweenU(bottom, top int) (result int) {
	// формула генерации случайных чисел по заданному диапазону
	// где bottom - минимальное число из желаемого диапазона
	// top - верхнаяя граница, ширина выборки
	rnd := int(RndU())
	div := rnd % top
	diff := top - div
	if diff > bottom {
		result = bottom + div
	} else {
		result = div
	}
	return
}

func main() {
	// shake the generator!
	SRnd64(time.Now().Unix())
	fmt.Println("Random:", RndBetweenU(1,9))	


	in := `first_name,last_name,username
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri`

	r := csv.NewReader(strings.NewReader(in))

	for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(record)
	}
	


}