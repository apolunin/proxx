# PROXX Game

This program implements a simple console version of the [Proxx](https://proxx.app/) game.
Currently flag functionality is not implemented.

## Installation and usage

This program requires Go 1.19+. It takes board size and number of holes on the board
as command line arguments. Use the following command to get detailed help:
```shell
go run ./cmd/proxx/proxx.go -h
```
The output would look similar to the following:
```
Usage of ./proxx:
  -k int
    	specifies number of holes (should be in range [5;9792]) (default 5)
  -n int
    	specifies game field size (should be in range [4;99]) (default 4)
```

## Gameplay

Run the game, specify board size and number of holes as needed.
For example the following command will run the game with 5x5 field having 8 holes:
```
go run ./cmd/proxx/proxx.go -n 5 -k 8
```
The program will display a board with rows and columns marked by numbers
from 1 to 5:
```
     1  2  3  4  5
  1  □  □  □  □  □
  2  □  □  □  □  □
  3  □  □  □  □  □
  4  □  □  □  □  □
  5  □  □  □  □  □

row, col = 
```
Below the board a prompt is displayed.
It asks a user to input row and column numbers, separated by comma.
When user inputs numbers and hits `return` key a game reveals the selected cell.
If it is a hole then game is over. User lost. A cell on a board may be displayed
via the following characters:
* □ - marks unrevealed cell (i.e. it's content is unknown to the user);
* ◯ - marks a hole;
* 0-8 - numbers 0-8 specify total number of holes in adjacent cells;
```

     1  2  3  4  5
  1  □  □  □  □  □
  2  □  □  □  □  □
  3  □  □  □  □  □
  4  □  □  □  □  □
  5  □  □  □  □  □

row, col = 1, 2

     1  2  3  4  5
  1  0  0  1  □  □
  2  2  2  2  □  □
  3  □  □  □  □  □
  4  □  □  □  □  □
  5  □  □  □  □  □

row, col = 
```
When a revealed cell has no holes in adjacent cells,
adjacent cells are revealed automatically as displayed above.

The game ends when user wins or loses.
A user loses when a cell containing a hole is revealed:
```
row, col = 1, 4

     1  2  3  4  5
  1  0  0  1  ◯  1
  2  2  2  2  2  2
  3  □  □  □  2  □
  4  □  □  □  □  □
  5  □  □  □  □  □

*** YOU LOST ***
```
A user wins when only hole cells are unrevealed:
```
row, col = 5,4

     1  2  3  4  5
  1  0  2  □  3  1
  2  1  3  □  3  □
  3  2  □  3  2  1
  4  3  □  4  2  1
  5  2  □  □  2  □

*** YOU WON ***
```
Happy playing!