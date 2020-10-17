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

	// var records [][]string

	// fmt.Println("Randomazing...")
	// test_random()

	// fmt.Println("Reading csv")
	// records = test_read_csv()

	// fmt.Println("Writing down...")
	// test_write_csv(records)

	// game_board := [][]string{
	// 	{"*", "|", "*", "|", "*"},
	// 	{"-", "|", "-", "|", "-"},
	// 	{"*", "|", "*", "|", "*"},
	// 	{"-", "|", "-", "|", "-"},
	// 	{"*", "|", "*", "|", "*"},
	// }

	// for _, line := range game_board {
	// 	fmt.Println(line)
	// }

	// game_board[2][2] = "X"
	// game_board[0][2] = "O"
	// fmt.Println()

	// for _, line := range game_board {
	// 	fmt.Println(line)
	// }

	// put_symbol(&game_board, 4, "X")
	// print_game_board(game_board)

	// f, err := os.OpenFile("game_board.csv", os.O_RDWR|os.O_CREATE, 0755)
	// defer func(err error) {
	// 	if err := f.Close(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }(err)

	// fmt.Println("File opened successfuly")

	// fmt.Println("Write to CSV")
	// wtire_game_board_to_csv("game_board.csv", empty_game_board())
	// fmt.Println("Read from CSV")
	// game_board := read_geme_board_from_cvs("game_board.csv")
	// print_game_board(game_board)

	// put_symbol(&game_board, 1, "O")
	// put_symbol(&game_board, 2, "O")
	// put_symbol(&game_board, 3, "O")

	// win_sym, err := match_three(game_board, []int{1, 2, 3})
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("Winner is palayer with with symbol:", win_sym, "!")
	// }

	// // test game ending
	// var winner string
	// game_end := test_game_end(game_board, &winner)
	// fmt.Println(game_end, winner)

	// // test multiwriting
	// put_symbol(&game_board, 7, "X")
	// put_symbol(&game_board, 8, "X")
	// put_symbol(&game_board, 9, "X")
	// print_game_board(game_board)
	// err = put_symbol(&game_board, 9, "X")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	test_gameplay()
}

func test_gameplay() {

	fmt.Println("Start the game!")

	// shake the generator!
	SRnd64(time.Now().Unix())

	game_board := read_geme_board_from_cvs("game_board.csv")
	turn := "X"

	print_game_board(game_board)

	var winner string

	for game_end := false; !game_end; game_end = test_game_end(game_board, &winner) {
		// try until success
		for stop := false; !stop; {
			pos := RndBetweenU(1, 10)
			err := put_symbol(&game_board, pos, turn)
			if err == nil {
				stop = true
			}
		}

		print_game_board(game_board)

		// end turn
		if turn == "X" {
			turn = "O"
		} else {
			turn = "X"
		}

		// Calling Sleep method
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Winner is a player with symbol", winner, "!")
}

func match_three(board [][]string, line []int) (string, error) {
	// we should to have exactly len(line)-1 matchings
	// otherwise it doesn't matching
	match_count := 0

	// get first value
	head := line[:1]
	rest := line[1:]
	cur_pos := head[0]

	// for all line get next value and compare with perv value
	for _, next_pos := range rest {
		// fmt.Println(cur_pos, ",", next_pos)

		// if it matches then incremetn match count
		cur_val, _ := get_at_pos(&board, cur_pos)
		next_val, _ := get_at_pos(&board, next_pos)

		// skip empty places
		if cur_val == "*" || next_val == "*" {
			continue
		}

		if cur_val == next_val {
			match_count++
		}

		// move further
		cur_pos = next_pos
	}

	// test is there a 'match three' or not
	if match_count == len(line)-1 {
		symbol_of_winner, _ := get_at_pos(&board, head[0])
		// fmt.Println("Match!", symbol_of_winner)
		return symbol_of_winner, nil
	} else {
		// fmt.Println("Not match!")
		return "", errors.New("Not match")
	}

}

func test_game_end(board [][]string, winner *string) (geme_end bool) {
	matchings := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{1, 5, 9},
		{3, 5, 7},
	}

	for _, line := range matchings {
		win_sym, err := match_three(board, line)
		if err == nil {
			*winner = win_sym
			return true
		}
	}

	return false
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
	defer func(err error) {
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
OPEN:
	f, err := os.Open(fullpath)
	if err != nil {
		// log.Fatal("Can't open file: ", fullpath)
		// create a new geme board if there is no any one
		wtire_game_board_to_csv(fullpath, empty_game_board())
		goto OPEN
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
	// fmt.Println("Lnes read:", line_count)

	// if there isn't a correct game board just to initialize it now by an empty board
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
	case 1:
		value = (*board)[0][0]
		break
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
	default: // out of range
		err = errors.New("Out of range")
		break
	}

	return
}

func set_at_pos(board *[][]string, pos int, symbol string) error {
	var err error = nil

	switch pos {
	case 1:
		(*board)[0][0] = symbol
		break
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
	default: // out of range
		err = errors.New("Out of range")
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
		// fmt.Println("pos", pos, "is busy by symbol", value)
		err = errors.New("Position is busy")
	} else if value == "*" {
		// fmt.Println("put symbol", symbol, "to pos", pos)
		err = set_at_pos(board, pos, symbol)
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
