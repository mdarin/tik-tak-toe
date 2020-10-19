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
	"net/http"
	"os"
	"strings"
	"time"
	// "net/url"
	"io/ioutil"
	// "bufio" // to scan and tokenize buffered input data from an io.Reader source
	"strconv"
	// "regexp"
	// "errors" // for errors.New()
	// "os" // for operations with dirs
	// "time" // for sleep
	// "sort" // for sorging
	// "encoding/base64"
	"bytes"
	"encoding/json"
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

	// test is the game in progress or it's a new gema

	// if new game then create new game board and lock

	// read game board from file

	// iterate the game as one turn for every player

	// test game end

	// if geme ended then remove game board and lock files
	// send the boad state and winner info
	// otherwise
	// send only board state and some catchword for supporting players

	var turn string = "X"
	var winner string
	var turns_count int
	var game_end bool = false

	// shake the generator!
	SRnd64(time.Now().Unix())

	// test is the game in progress or it's a new gema
	state, err := read_lock_from_csv("lock.csv")
	// if new game then create new game board and lock
	if err != nil {
		// create a new geme board and lock
		wtire_lock_to_csv("lock.csv", default_state())
		wtire_game_board_to_csv("game_board.csv", empty_game_board())
		state = default_state()
		s := state[0][0]
		turns_count, _ = strconv.Atoi(s)
		fmt.Println("NO LOCK New game:", turns_count)
	} else {
		s := state[0][0]
		turns_count, _ = strconv.Atoi(s)
		fmt.Println("LOCKED Continue game:", turns_count)
	}

	// read game board from file
	game_board := read_geme_board_from_cvs("game_board.csv")

	// iterate the game as one turn for every player
	for i := 0; i < 2; i++ {
		// fmt.Println(turn)
		// try until success
		for stop := false; !stop && turns_count > 0; {
			pos := RndBetweenU(1, 10)
			// fmt.Println("pos:", pos, " turn:", turn, " remains turns:", turns_count)
			err := put_symbol(&game_board, pos, turn)
			if err == nil {
				stop = true
			}
		}

		// end turn
		if turn == "X" {
			turn = "O"
		} else {
			turn = "X"
		}

		turns_count--
	}

	// save and show game board
	wtire_game_board_to_csv("game_board.csv", game_board)
	// print_game_board(game_board)
	fmt.Printf(to_string(game_board))

	// test game end
	game_end = test_game_end(game_board, &winner)

	// save state into the lock file
	s := strconv.Itoa(turns_count)
	state[0][0] = s
	wtire_lock_to_csv("lock.csv", state)

	// if geme ended then remove game board and lock files
	// send the boad state and winner info
	if game_end || turns_count < 1 {
		// unlock
		err := os.Remove("lock.csv")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("UNLOCKED")

		// send state actualization
		if winner == "X" || winner == "O" {
			fmt.Println("Winner is a player with symbol", winner)
			// send("Winner is a player with symbol " + winner)
		} else {
			fmt.Println("It is a drawn game")
			// send("It is a drawn game")
		}
	}

	// otherwise
	// send only board state and some catchword for supporting players
	message := "```" + to_string(game_board) + "```"
	send(message)

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

	//test_gameplay()

	// // test to string
	// e := empty_game_board()
	// put_symbol(&e, 1, "X")
	// s := to_string(e)
	// fmt.Printf(s)

	// // test locking

	// for i := 0; i < 2; i++ {
	// 	state, err := read_lock_from_csv("lock.csv")
	// 	if err != nil {
	// 		wtire_lock_to_csv("lock.csv", default_state())
	// 		state, _ = read_lock_from_csv("lock.csv")
	// 		s := state[0][0]
	// 		turns_count, _ := strconv.Atoi(s)
	// 		s = strconv.Itoa(turns_count-1)
	// 		state[0][0] = s
	// 		wtire_lock_to_csv("lock.csv", state)
	// 		fmt.Println("NO LOCK New game:", turns_count)
	// 	} else {
	// 		turns_count := state[0][0]
	// 		fmt.Println("LOCKED Continue game:", turns_count)
	// 	}
	// 	if i == 1 {
	// 		err := os.Remove("lock.csv")
	// 		if err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println("UNLOCKED")
	// 	}
	// }
} // eof main

func default_state() [][]string {
	return [][]string{
		{"9"},
	}
}

func test_gameplay() {

	fmt.Println("Start the game!")

	// shake the generator!
	SRnd64(time.Now().Unix())

	// create a new geme board if there is no any one
	wtire_game_board_to_csv("game_board.csv", empty_game_board())

	game_board := read_geme_board_from_cvs("game_board.csv")
	turn := "X"

	print_game_board(game_board)

	var winner string
	turns_count := 9
	for game_end := false; !game_end && turns_count > 0; {
		fmt.Println()

		// try until success
		for stop := false; !stop; {
			pos := RndBetweenU(1, 10)
			fmt.Println("pos:", pos, " turn:", turn, " remains turns:", turns_count)
			err := put_symbol(&game_board, pos, turn)
			if err == nil {
				stop = true
			}
		}

		print_game_board(game_board)
		// wtire_game_board_to_csv("game_board.csv", game_board)

		// end turn
		if turn == "X" {
			turn = "O"
		} else {
			turn = "X"
		}

		game_end = test_game_end(game_board, &winner)
		turns_count--

		// Calling Sleep method
		// time.Sleep(1 * time.Second)
	}

	if winner == "X" || winner == "O" {
		fmt.Println("Winner is a player with symbol", winner)
		send("Winner is a player with symbol " + winner)
	} else {
		fmt.Println("It is a drawn game")
		send("It is a drawn game")
	}

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
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
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
		{"-", "+", "-", "+", "-"},
		{"*", "|", "*", "|", "*"},
		{"-", "+", "-", "+", "-"},
		{"*", "|", "*", "|", "*"},
	}
}

func wtire_lock_to_csv(fullpath string, state [][]string) error {
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

	for _, record := range state {
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

func read_lock_from_csv(fullpath string) ([][]string, error) {
	var records [][]string = nil

	f, err := os.Open(fullpath)
	if err != nil {
		// log.Fatal("Can't open lock file: ", fullpath)
		return records, err
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
		line_count++
	}

	return records, err
}

func read_geme_board_from_cvs(fullpath string) (records [][]string) {
	f, err := os.Open(fullpath)
	if err != nil {
		log.Fatal("Can't open file: ", fullpath)
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
		// fmt.Println(record)
		line_count++
	}
	// fmt.Println("Lnes read:", line_count)

	// if there isn't a correct game board just to initialize it now by an empty board
	if line_count != 5 {
		wtire_game_board_to_csv(fullpath, empty_game_board())
	}

	return records
}

func to_string(board [][]string) string {
	result := ""

	for _, line := range board {
		for _, sym := range line {
			if sym == "*" {
				result += " "
			} else {
				result += sym
			}
		}
		result += "\n"
	}

	return result
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

// {
// 	"username": "Dungeones Keeperra",
// 	"icon_emoji": ":smiling_imp:",
// 	"channel": "#backend",
// 	"text": "День расплаты настал!:fire:\nСегодня необходимо уделить время ревью.\nПерейдите в канал <https://aeoner.slack.com/archives/CNHR5NV45|#mr_backend> и доведите дело до конца.\nИсполните свой долг самурая."
// }
// Creatе an issue
func send(text string) []byte {
	// prepare body payload
	payload, _ := json.Marshal(struct {
		Username   string `json:"username"`
		Icon_emoji string `json:"icon_emoji"`
		Channel    string `json:"channel"`
		Text       string `json:"text"`
	}{
		Username:   "Dirty J",
		Icon_emoji: ":ghost:",
		Channel:    "#test_ch",
		Text:       text,
	})

	fmt.Println(string(payload))

	// prepare request
	//TODO: user path methods to concatinate!
	url := "https://hooks.slack.com/services/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"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	// req.SetBasicAuth(AEON_USER, AEON_USER_TOKEN)

	fmt.Println(req)

	// prepare http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}

	}()

	fmt.Println(resp.StatusCode)

	// prpcess response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return body
}
