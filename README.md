# README

Simple practice session for Go.
I created a Tris Game that create the whole tree of moves and solutions.
The aim is to create a Computer player that always win or drew

# Run

```
go run pkg/main.go
```

Use number from 0 to 8 to place your move

# Enhancement

The game has complete knowledge but bad strategy: the method `func (node *Node) BestMove() Bits` need to be improved
