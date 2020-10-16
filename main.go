/***
 *
 *
 * csv - https://golang-blog.blogspot.com/2020/06/csv-package-in-golang.html
 **/
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	_ "math"
	"os"
	"strings"
	"time"
)

var (
	gRndSeed uint32  = 1 // последнее случайное число
	was      int     = 0 // была ли вычислена пара чисел
	r        float64 = 0 // предыдущее число
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
	gRndSeed = gRndSeed*uint32(1664525) + uint32(1013904223)
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

	var records [][]string

	fmt.Println("Randomazing...")
	test_random()

	fmt.Println("Reading csv")
	records = test_read_csv()

	fmt.Println("Writing down...")
	test_write_csv(records)

	game_board := [][]string{
		{"*", "|", "*", "|", "*"},
		{"-", "|", "-", "|", "-"},
		{"*", "|", "*", "|", "*"},
		{"-", "|", "-", "|", "-"},
		{"*", "|", "*", "|", "*"},
	}

	for _, line := range game_board {
		fmt.Println(line)
	}

	game_board[2][2] = "X"
	game_board[0][2] = "O"
	fmt.Println()

	for _, line := range game_board {
		fmt.Println(line)
	}

	put_symbol(&game_board, 4, "X")
	print_game_board(game_board)

	f, err := os.OpenFile("game_board.csv", os.O_RDWR|os.O_CREATE, 0755)
	defer func(err error) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}(err)

	fmt.Println("File opened successfuly")

	fmt.Println("Write to CSV")
	// wtire_game_board_to_csv("game_board.csv", game_board)
	fmt.Println("Read from CSV")
	game_board = read_geme_board_from_cvs("game_board.csv")
	print_game_board(game_board)


}

func empty_game_board() [][]string {
	return [][]string{
		{"*", "|", "*", "|", "*"},
		{"-", "|", "-", "|", "-"},
		{"*", "|", "*", "|", "*"},
		{"-", "|", "-", "|", "-"},
		{"*", "|", "*", "|", "*"},
	}
}

func wtire_game_board_to_csv(fullpath string, board [][]string) error {
	f, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE, 0755)
	defer func (err error) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}(err)
	if err != nil {
		log.Fatal("Can't create new file", fullpath)	
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, record := range board {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Записываем любые буферизованные данные в подлежащий writer (стандартный вывод).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	return err
}

func read_geme_board_from_cvs(fullpath string) (records [][]string) {
	f, err := os.Open(fullpath)
	if err != nil {
		log.Fatal("Can't open file", fullpath)

	}
	defer func(err error) {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}(err)

	r := csv.NewReader(f)

	line_count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, record)
		fmt.Println(record)
		line_count++
	}
	fmt.Println("Lnes read:", line_count)
	
	// if there isn't a correct game board just to initialize it now by empty board
	if line_count != 5 {
		wtire_game_board_to_csv(fullpath, empty_game_board())
	}

	return records
}
func print_game_board(board [][]string) {
	for _, line := range board {
		fmt.Println(line)
	}
}

func get_at_pos(board *[][]string, pos int) (value string, err error) {
	err = nil
	switch pos {
	case 2:
		value = (*board)[0][2]
		break
	case 3:
		value = (*board)[0][4]
		break
	case 4:
		value = (*board)[2][0]
		break
	case 5:
		value = (*board)[2][2]
		break
	case 6:
		value = (*board)[2][4]
		break
	case 7:
		value = (*board)[4][0]
		break
	case 8:
		value = (*board)[4][2]
		break
	case 9:
		value = (*board)[4][4]
		break
	default: // 1st position
		value = (*board)[0][0]
		break
	}

	return
}

func set_at_pos(board *[][]string, pos int, symbol string) error {
	var err error = nil

	switch pos {
	case 2:
		(*board)[0][2] = symbol
		break
	case 3:
		(*board)[0][4] = symbol
		break
	case 4:
		(*board)[2][0] = symbol
		break
	case 5:
		(*board)[2][2] = symbol
		break
	case 6:
		(*board)[2][4] = symbol
		break
	case 7:
		(*board)[4][0] = symbol
		break
	case 8:
		(*board)[4][2] = symbol
		break
	case 9:
		(*board)[4][4] = symbol
		break
	default: // 1st position
		(*board)[0][0] = symbol
		break
	}

	return err
}

func put_symbol(board *[][]string, pos int, symbol string) error {
	var err error = nil

	if symbol != "X" && symbol != "O" {
		return errors.New("Incorrect symbol")
	}

	value, _ := get_at_pos(board, pos)
	if value == "X" || value == "O" {
		fmt.Println("pos", pos, "is busy by symbol", value)
		err = errors.New("Position is busy")
	} else if value == "*" {
		fmt.Println("put symbol", symbol, "to pos", pos)
		set_at_pos(board, pos, symbol)
	}

	return err
}

func test_random() {
	// shake the generator!
	SRnd64(time.Now().Unix())
	// use generator
	fmt.Println("Random:", RndBetweenU(1, 9))
}

func test_read_csv() (records [][]string) {
	in := `first_name,last_name,username
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri`

	r := csv.NewReader(strings.NewReader(in))

	lineno := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, record)
		fmt.Println(record)
		lineno++
	}
	fmt.Println("Lnes read:", lineno)

	return records
}

func test_write_csv(records [][]string) {
	w := csv.NewWriter(os.Stdout)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Записываем любые буферизованные данные в подлежащий writer (стандартный вывод).
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
