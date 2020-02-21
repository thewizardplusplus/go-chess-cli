# go-chess-cli

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-chess-cli?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-chess-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-chess-cli)](https://goreportcard.com/report/github.com/thewizardplusplus/go-chess-cli)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-chess-cli.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-chess-cli)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-chess-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-chess-cli)

The chess program with a terminal-based interface.

_**Disclaimer:** this program was written directly on an Android smartphone with the AnGoIde IDE._

## Features

- displaying a board:
  - by symbols (to choose):
    - ASCII;
    - Unicode;
  - by size (to choose):
    - terse;
    - wide;
  - by colors (to choose):
    - monochrome;
    - colorful;
  - misc.:
    - marking searching process;
    - placing a human side at bottom;
- interacting via text commands (moves in [pure algebraic coordinate notation](https://www.chessprogramming.org/Algebraic_Chess_Notation#Pure_coordinate_notation));
- options:
  - initial position in [Forsyth–Edwards notation](https://en.wikipedia.org/wiki/Forsyth–Edwards_Notation);
  - human color (i.e. a computer can move first):
    - support automatic random selecting (optional);
  - move searching restrictions:
    - maximal size of the [transposition table](https://www.chessprogramming.org/Transposition_Table);
    - deep of move searching;
    - duration of move searching;
  - displaying:
    - switching between ASCII/Unicode modes;
    - switching between terse/wide modes;
    - switching between monochrome/colorful modes.

## Installation

```
$ go get github.com/thewizardplusplus/go-chess-cli/...
```

## Usage

```
$ go-chess-cli -h | -help | --help
$ go-chess-cli [options]
```

Options:

- `-h`, `-help`, `--help` &mdash; show the help message and exit;
- `-cacheSize ITEMS` &mdash; maximal cache size (default: `1000000`, i.e. one million);
- `-colorfulBoard {false|true}` &mdash; use colors to display the board (default: `true`; for inverting use `-colorfulBoard=false`);
- `-colorfulPieces {false|true}` &mdash; use colors to display pieces (default: `true`; for inverting use `-colorfulPieces=false`);
- `-deep INTEGER` &mdash; search deep (default: `5`);
- `-duration DURATION` &mdash; search duration (e.g. `72h3m0.5s`; default: `5s`);
- `-fen STRING` &mdash; board in FEN (default: `rnbqk/ppppp/5/PPPPP/RNBQK`, i.e. Gardner's minichess);
- `-humanColor {random|black|white}` &mdash; human color (default: `random`);
- `-pieceBlackColor INTEGER` &mdash; SGR parameter for ANSI escape sequences for setting a color of black pieces (default: `34`; see for details: https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit);
- `-pieceWhiteColor INTEGER` &mdash; SGR parameter for ANSI escape sequences for setting a color of white pieces (default: `31`; see for details: https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit);
- `-squareBlackColor INTEGER` &mdash; SGR parameter for ANSI escape sequences for setting a color of black squares (default: `40`; see for details: https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit);
- `-squareWhiteColor INTEGER` &mdash; SGR parameter for ANSI escape sequences for setting a color of white squares (default: `47`; see for details: https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit);
- `-unicode {false|true}` &mdash; use Unicode to display pieces (default: `true`; for inverting use `-unicode=false`);
- `-wide {false|true}` &mdash; display the board wide (default: `true`; for inverting use `-wide=false`).

## Examples

`ascii.DecodeColor()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
)

func main() {
	color, _ := ascii.DecodeColor("white")
	fmt.Printf("%v\n", color)

	// Output: 1
}
```

`ascii.EncodeColor()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-chess-models"
)

func main() {
	color := ascii.EncodeColor(models.White)
	fmt.Printf("%v\n", color)

	// Output: white
}
```

`ascii.PieceStorageEncoder.EncodePieceStorage()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
	const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
	storage, _ := uci.DecodePieceStorage(fen, pieces.NewPiece, models.NewBoard)
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		ascii.Margins{},
		ascii.WithoutColor,
		models.Black,
		1,
	)
	fmt.Printf("%v\n", encoder.EncodePieceStorage(storage))

	// Output:
	// 8rxxxkxxr
	// 7pxppqpbx
	// 6bnxxpnpx
	// 5xxxPNxxx
	// 4xpxxPxxx
	// 3xxNxxQxp
	// 2PPPBBPPP
	// 1RxxxKxxR
	//  abcdefgh
}
```

`ascii.PieceStorageEncoder.EncodePieceStorage()` with margins:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
	margins := ascii.Margins{
		Piece: ascii.PieceMargins{
			HorizontalMargins: ascii.HorizontalMargins{
				Left: 1,
			},
		},
		Legend: ascii.LegendMargins{
			Rank: ascii.HorizontalMargins{
				Right: 1,
			},
		},
	}

	const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
	storage, _ := uci.DecodePieceStorage(fen, pieces.NewPiece, models.NewBoard)
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		margins,
		ascii.WithoutColor,
		models.Black,
		1,
	)
	fmt.Printf("%v\n", encoder.EncodePieceStorage(storage))

	// Output:
	// 8  r x x x k x x r
	// 7  p x p p q p b x
	// 6  b n x x p n p x
	// 5  x x x P N x x x
	// 4  x p x x P x x x
	// 3  x x N x x Q x p
	// 2  P P P B B P P P
	// 1  R x x x K x x R
	//    a b c d e f g h
}
```

`ascii.PieceStorageEncoder.EncodePieceStorage()` with colors:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/ascii"
	climodels "github.com/thewizardplusplus/go-chess-cli/models"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/encoding/uci"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
	colorizer := func(text string, color climodels.OptionalColor) string {
		var colorMark byte
		if color.IsSet {
			colorMark = ascii.EncodeColor(color.Value)[0]
		} else {
			colorMark = 'n'
		}

		return fmt.Sprintf("(%c%s)", colorMark, text)
	}

	const fen = "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R"
	storage, _ := uci.DecodePieceStorage(fen, pieces.NewPiece, models.NewBoard)
	encoder := ascii.NewPieceStorageEncoder(
		uci.EncodePiece,
		"x",
		ascii.Margins{},
		colorizer,
		models.Black,
		1,
	)
	fmt.Printf("%v\n", encoder.EncodePieceStorage(storage))

	// Output:
	// (n8)(wr)(bx)(wx)(bx)(wk)(bx)(wx)(br)
	// (n7)(bp)(wx)(bp)(wp)(bq)(wp)(bb)(wx)
	// (n6)(wb)(bn)(wx)(bx)(wp)(bn)(wp)(bx)
	// (n5)(bx)(wx)(bx)(wP)(bN)(wx)(bx)(wx)
	// (n4)(wx)(bp)(wx)(bx)(wP)(bx)(wx)(bx)
	// (n3)(bx)(wx)(bN)(wx)(bx)(wQ)(bx)(wp)
	// (n2)(wP)(bP)(wP)(bB)(wB)(bP)(wP)(bP)
	// (n1)(bR)(wx)(bx)(wx)(bK)(wx)(bx)(wR)
	// (n )(na)(nb)(nc)(nd)(ne)(nf)(ng)(nh)
}
```

`unicode.EncodePiece()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-chess-cli/encoding/unicode"
	models "github.com/thewizardplusplus/go-chess-models"
	"github.com/thewizardplusplus/go-chess-models/pieces"
)

func main() {
	fen := unicode.EncodePiece(pieces.NewBishop(models.White, models.Position{}))
	fmt.Printf("%v\n", fen)

	// Output: ♗
}
```

## License

The MIT License (MIT)

Copyright &copy; 2019-2020 thewizardplusplus
