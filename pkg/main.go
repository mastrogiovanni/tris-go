package main

import (
	"errors"
	"fmt"
)

type Bits uint16

const (
	P00 Bits = 1 << iota
	P10
	P20
	P01
	P11
	P21
	P02
	P12
	P22
)

const (
	NO_WINNER_YET uint8 = iota
	WIN_ODD             // Odd Height
	WIN_EVEN            // Even Height
	NO_WINNER           // Nobody
)

var solutions []Bits = []Bits{
	P00 | P11 | P22, // left top to right bottom
	P20 | P11 | P02, // left bottom to right rop
	P00 | P10 | P20, // horizontal, 1 row
	P01 | P11 | P21, // horizontal, 2 row
	P02 | P12 | P22, // horizontal, 3 row
	P00 | P01 | P02, // vertical, 1 col
	P10 | P11 | P12, // vertical, 2 col
	P20 | P21 | P22, // vertical, 3 col
}

const Full Bits = P00 | P01 | P02 | P10 | P11 | P12 | P20 | P21 | P22

type Node struct {
	Parent   *Node
	Children []*Node
	Move     Bits
	Height   uint8
	Remain   Bits
	Winner   uint8
}

type CountWinners struct {
	NoWinner   uint
	OddWinner  uint
	EvenWinner uint
}

func newNode(parent *Node, move Bits) *Node {
	if parent == nil {
		return &Node{
			Parent:   parent,
			Children: []*Node{},
			Move:     move,
			Height:   0,
			Remain:   Bits(0),
			Winner:   NO_WINNER_YET,
		}
	} else {
		return &Node{
			Parent:   parent,
			Children: []*Node{},
			Move:     move,
			Height:   parent.Height + 1,
			Remain:   parent.Remain | move,
			Winner:   NO_WINNER_YET,
		}
	}
}

func IndexToMove(index int) Bits {
	switch index {
	case 0:
		return P00
	case 1:
		return P10
	case 2:
		return P20
	case 3:
		return P01
	case 4:
		return P11
	case 5:
		return P21
	case 6:
		return P02
	case 7:
		return P12
	case 8:
		return P22
	}
	return 0
}

func MoveToIndex(move Bits) uint8 {
	if move&P00 != 0 {
		return 0
	}
	if move&P10 != 0 {
		return 1
	}
	if move&P20 != 0 {
		return 2
	}
	if move&P01 != 0 {
		return 3
	}
	if move&P11 != 0 {
		return 4
	}
	if move&P21 != 0 {
		return 5
	}
	if move&P02 != 0 {
		return 6
	}
	if move&P12 != 0 {
		return 7
	}
	if move&P22 != 0 {
		return 8
	}
	return 255
}

func (node *Node) Print(numbers bool) {
	var moves []rune
	if numbers {
		moves = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8'}
	} else {
		moves = []rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	}
	c := rune('x')
	for node.Height > 0 {
		index := MoveToIndex(node.Move)
		if node.Height%2 == 0 {
			c = rune('o')
		} else {
			c = rune('x')
		}
		moves[index] = c
		node = node.Parent
	}
	fmt.Printf("%c│%c│%c\n", moves[0], moves[1], moves[2])
	fmt.Println("─┼─┼─")
	fmt.Printf("%c│%c│%c\n", moves[3], moves[4], moves[5])
	fmt.Println("─┼─┼─")
	fmt.Printf("%c│%c│%c\n\n", moves[6], moves[7], moves[8])
}

func (node *Node) IsFull() bool {
	return node.Remain == Full
}

func (node *Node) Check() bool {
	bitmap := Bits(0)
	parity := node.Height % 2
	for node.Height > 0 {
		if node.Height%2 == parity {
			bitmap = bitmap | node.Move
		}
		node = node.Parent
	}

	for s := range solutions {
		if bitmap&solutions[s] == solutions[s] {
			return true
		}
	}
	return false
}

func (node *Node) Set(move Bits) (*Node, error) {

	if node.Winner != NO_WINNER_YET {
		return nil, errors.New("the game is over")
	}

	if node.Remain&move != 0 {
		return nil, errors.New("move not available")
	}

	for i := range node.Children {
		if node.Children[i].Move == move {
			return node.Children[i], nil
		}
	}

	newNode := newNode(node, move)
	node.Children = append(node.Children, newNode)

	if newNode.Check() {
		if newNode.Height%2 == 0 {
			newNode.Winner = WIN_EVEN
		} else {
			newNode.Winner = WIN_ODD
		}
	} else {
		if newNode.IsFull() {
			newNode.Winner = NO_WINNER
		}
	}

	return newNode, nil
}

func Explore(node *Node) {
	for i := 0; i < 9; i++ {
		move := IndexToMove(i)
		newNode, err := node.Set(move)
		if err == nil {
			if newNode.Winner == NO_WINNER_YET {
				Explore(newNode)
			}
		}
	}
}

type Visitor func(interface{}, *Node)

func (root *Node) DepthVisit(context interface{}, visit Visitor) {
	if root != nil {
		visit(context, root)
		for i := range root.Children {
			root.Children[i].DepthVisit(context, visit)
		}
	}
}

func (node *Node) PrintWinners() {
	node.DepthVisit(nil, func(context interface{}, node *Node) {
		if node.Winner == WIN_EVEN {
			fmt.Println("O win")
			node.Print(false)
		}
		if node.Winner == WIN_ODD {
			fmt.Println("X win")
			node.Print(false)
		}
		if node.Winner == NO_WINNER {
			fmt.Println("Nobody win")
			node.Print(false)
		}
	})
}

func (node *Node) CountWinners() *CountWinners {
	context := &CountWinners{
		NoWinner:   0,
		OddWinner:  0,
		EvenWinner: 0,
	}
	node.DepthVisit(context, func(context interface{}, node *Node) {
		var c *CountWinners = context.(*CountWinners)
		if node.Winner == WIN_EVEN {
			c.EvenWinner = c.EvenWinner + 1
		}
		if node.Winner == WIN_ODD {
			c.OddWinner = c.OddWinner + 1
		}
		if node.Winner == NO_WINNER {
			c.NoWinner = c.NoWinner + 1
		}
	})
	return context
}

func (node *Node) BestMove() Bits {
	var max float64 = -1
	var maxI int = 0
	for i := range node.Children {
		winners := node.Children[i].CountWinners()
		// fmt.Println("Move:", MoveToIndex(node.Children[i].Move))
		if node.Children[i].Height%2 == 0 {
			probability := float64(winners.EvenWinner+winners.NoWinner) / float64(winners.EvenWinner+winners.OddWinner+winners.NoWinner)
			// fmt.Println(probability)
			if probability > max {
				max = probability
				maxI = i
			}
		} else {
			probability := float64(winners.OddWinner+winners.NoWinner) / float64(winners.EvenWinner+winners.OddWinner+winners.NoWinner)
			// fmt.Println(probability)
			if probability > max {
				max = probability
				maxI = i
			}
		}
	}
	// fmt.Println("Best choice", MoveToIndex(node.Children[maxI].Move))
	return node.Children[maxI].Move
}

func (node *Node) Count() int {
	tot := 1
	for i := range node.Children {
		tot = tot + node.Children[i].Count()
	}
	return tot
}

func IsEnded(node *Node, who string) bool {
	isTris := node.Check()
	if isTris {
		fmt.Println(who, " Wins")
		return true
	}
	if node.IsFull() {
		fmt.Println("Nobody Wins")
		return true
	}
	return false
}

func main() {

	root := newNode(nil, Bits(0))

	Explore(root)

	playerStart := true

	for {

		fmt.Println("----------- New Game -------------")

		game := root

		if playerStart {

			game.Print(true)

			for {
				fmt.Print("Player move: ")
				var index int
				for {
					parsed, err := fmt.Scanf("%d", &index)
					if parsed == 1 {
						fmt.Println(index)
						game, _ = game.Set(IndexToMove(index))
						game.Print(false)
						break
					} else {
						fmt.Println(err)
					}
				}
				if IsEnded(game, "Player") {
					break
				}

				fmt.Println("Computer move: ")
				move := game.BestMove()
				game, _ = game.Set(move)
				game.Print(true)
				if IsEnded(game, "Computer") {
					break
				}
			}

		} else {

			game.Print(false)

			for {
				fmt.Println("Computer move: ")
				move := game.BestMove()
				game, _ = game.Set(move)
				game.Print(true)
				if IsEnded(game, "Computer") {
					break
				}

				fmt.Print("Player move: ")
				var index int
				for {
					parsed, err := fmt.Scanf("%d", &index)
					if parsed == 1 {
						fmt.Println(index)
						game, _ = game.Set(IndexToMove(index))
						game.Print(false)
						break
					} else {
						fmt.Println(err)
					}
				}
				if IsEnded(game, "Player") {
					break
				}

			}

		}

		playerStart = !playerStart

	}

}
