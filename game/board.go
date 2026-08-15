package main

import (
	"fmt"
	"unicode"
)

const startPositionFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR" // 9+9+

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
				file += int(c - '0') // start counting from '0' in ascii till c to get int conversion
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

func main() {
	nranks := 8
	nfiles := 8

	board, err := FENToBoard(startPositionFEN, nranks, nfiles)
	if err != nil {
		fmt.Println(err)
		return
	}

	board.print()

}
