package main

import (
	"fmt"
	"unicode"
)

const startPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"

type PieceType byte

const (
	None   PieceType = '·'
	Pawn   PieceType = 'p'
	Knight PieceType = 'n'
	Bishop PieceType = 'b'
	Rook   PieceType = 'r'
	Queen  PieceType = 'q'
	King   PieceType = 'k'
)

type Color byte

const (
	White Color = 'w'
	Black Color = 'b'
)

type Piece struct {
	Piece PieceType
	Color Color
}

type Board [][]Piece

var colorToString = map[Color]string{
	White: "white",
	Black: "black",
}

func FENToBoard(FEN string, nranks, nfiles int) (Board, error) {
	rank := 0
	file := 0

	board := make([][]Piece, nranks)
	for i := range board {
		board[i] = make([]Piece, nfiles)

		for j := range board {
			board[i][j] = Piece{Piece: None}
		}
	}

	for i, c := range FEN {
		switch c {
		case '/':
			if file != nfiles {
				return nil, fmt.Errorf("Invalid FEN at index %d: rank contains %d files instead of %d",
					i, file, nfiles)
			}
			rank++
			file = 0

		case 'p', 'n', 'b', 'r', 'q', 'k', 'P', 'N', 'B', 'R', 'Q', 'K':
			color := White
			if c >= 'a' && c <= 'z' {
				color = Black
			}

			// if uppercase convert to lowercase bc needed for conversion to PieceType
			if c >= 'A' && c <= 'Z' {
				c = unicode.ToLower(c)
			}

			board[rank][file] = Piece{Piece: PieceType(c), Color: color}
			file++

		default:
			if c >= '1' && c <= '8' {
				file += int(c - '0') // counts how far away c is from '0' in ascii table
			} else {
				return nil, fmt.Errorf("Invalid FEN at index %d: character '%c' is not supported",
					i, c)
			}
		}
	}

	return board, nil
}

func (b Board) print() {
	notations := "87654321"
	str := "  ┌────────────────────────┐\n"

	for rankIdx, rank := range b {

		str += string(notations[rankIdx])
		str += " │"
		for _, piece := range rank {
			r := rune(piece.Piece)

			if piece.Color == White {
				r = unicode.ToUpper(r)
			}

			str = fmt.Sprintf("%s %c ", str, r)

		}
		str += "│\n"
	}
	str += "  └────────────────────────┘\n"
	str += "    a  b  c  d  e  f  g  h\n"

	fmt.Print(str)
}

func (b Board) move(from, to, turn string) error {
	nFiles := len(b[0])
	nRanks := len(b)

	fromFile := int(from[0] - 'a')            // counts how far away char is from 'a'/'1' in ascii table
	fromRank := nRanks - 1 - int(from[1]-'1') // board starts top left but notation bottom left

	toFile := int(to[0] - 'a')
	toRank := nRanks - 1 - int(to[1]-'1')

	if fromFile < 0 || fromFile >= nFiles ||
		toFile < 0 || toFile >= nFiles ||
		fromRank < 0 || fromRank >= nRanks ||
		toRank < 0 || toRank >= nRanks {
		return fmt.Errorf("error: move (%s %s) is out of bounds", from, to)
	}

	fromPiece := b[fromRank][fromFile]
	if colorToString[fromPiece.Color] != turn {
		return fmt.Errorf("error: %s cannot move a piece that isn't theirs", colorToString[fromPiece.Color])
	}

	b[toRank][toFile] = b[fromRank][fromFile]
	b[fromRank][fromFile].Piece = None

	return nil
}

func (b Board) game() {
	turn := "white"
	for {
		var from string
		var to string

		b.print()

		fmt.Printf("%s move: ", turn)
		fmt.Scan(&from, &to)

		if len(from) != 2 || len(to) != 2 {
			fmt.Printf("error: incorrect notation. should be sth like `e2 e4` instead of `%s %s`\n", from, to)
			continue
		}

		err := b.move(from, to, turn)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if turn == "white" {
			turn = "black"
		} else {
			turn = "white"
		}
	}
}

func main() {
	nRanks := 8
	nFiles := 8

	board, err := FENToBoard(startPositionFEN, nRanks, nFiles)
	if err != nil {
		fmt.Println(err)
		return
	}

	board.game()
}
