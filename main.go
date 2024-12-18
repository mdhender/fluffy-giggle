// Copyright (c) 2024 Michael D Henderson. All rights reserved.

package main

import "fmt"

// a terminal "x" is implemented as
//   x := &Node{Sym:&Symbol{Terminal: true, Tsym: 'x'}}
// a non-terminal "A" is implemented as
//   a := &Header{Sym:'A'}
// a sequence of nodes S1, S2, S3 is implemented as
//   s3 := &Node{}
//   s2 := &Node{Suc: s3}
//   s1 := &Node{Suc: s2}
// a choice of nodes (a/k/a "alternates") S1, S2, S3 is implemented as
//   s3 := &Node{}
//   s2 := &Node{Alt: s3}
//   s1 := &Node{Alt: s2}
// a loop of node (a/k/a "zero or more") S is implemented as
//   empty := &Node{}
//   s := &Node{Alt: empty}
//   s.Suc = s

func main() {
	start := example5()
	fmt.Printf("example 5: %+v\n", start)
	firstSymbols := start.First(map[*Node]bool{})
	fmt.Printf("first %+v\n", firstSymbols)

	start = example5reduced()
	fmt.Printf("example 5 reduced: %+v\n", start)
	firstSymbols = start.First(map[*Node]bool{})
	fmt.Printf("first %+v\n", firstSymbols)
}

func example5() *Header {
	// Example 5 grammar
	//   A ::= 'x' | '(' B ')'
	//   B ::= A C
	//   C ::= { '+' A }

	empty := &Node{Terminal: true, Tsym: `ε`}
	lParen := &Node{Terminal: true, Tsym: `(`}
	plus := &Node{Terminal: true, Tsym: `+`}
	rParen := &Node{Terminal: true, Tsym: `)`}
	x := &Node{Terminal: true, Tsym: `x`}

	A, B, C := &Header{Sym: `A`}, &Header{Sym: `B`}, &Header{Sym: `C`}
	A.Entry = x
	x.Alt = lParen
	lParen.Suc = &Node{
		Nsym: A,
		Suc:  rParen,
	}
	B.Entry = &Node{
		Nsym: A,
		Suc:  &Node{Nsym: C},
	}
	C.Entry = plus
	plus.Alt = empty
	plus.Suc = &Node{Nsym: A, Suc: plus}

	// return A as the start symbol
	return A
}

func example5reduced() *Header {
	// Example 5 grammar (reduced)
	// A ::= 'x' | '(' A { '+' A } ')'
	A := &Header{Sym: `A`}
	x := &Node{Terminal: true, Tsym: `x`}
	lParen := &Node{Terminal: true, Tsym: `(`}
	empty := &Node{Terminal: true, Tsym: `ε`}
	plus := &Node{Terminal: true, Tsym: `+`}
	rParen := &Node{Terminal: true, Tsym: `)`}
	A.Entry = lParen
	lParen.Alt = x
	lParen.Suc = &Node{Nsym: A, Suc: plus}
	plus.Alt = empty
	plus.Suc = &Node{Nsym: A, Suc: plus}
	empty.Suc = rParen

	// return A as the start symbol
	return A
}

/*
 * type
 *   pointer = ^node;
 *   node =
 *		record suc, alt: pointer;
 *			case terminal: boolean of
 *				true: (tsym: char);
 *				false: (nsym: hpointer);
 *		end;
 */

type Node struct {
	// symbol for this node
	Terminal bool    // Discriminating flag: true for terminal, false for non-terminal
	Tsym     string  // Terminal symbol (valid if Terminal is true)
	Nsym     *Header // Non-terminal symbol (valid if Terminal is false)

	Suc *Node // Successor link
	Alt *Node // Alternative
}

/*
 * type
 * 	  hpointer = ^header;
 * 	  header =
 * 		record
 * 			entry: pointer;  // A pointer to a 'node' record
 * 			sym: char        // A character symbol
 * 		end;
 */

type Header struct {
	Sym   string // A character symbol
	Entry *Node  // Pointer to another Node
}

func (h *Header) First(visited map[*Node]bool) map[string]bool {
	var firstSymbols = map[string]bool{} // Set of terminal symbols
	for sym := range h.Entry.First(visited) {
		firstSymbols[sym] = true
	}
	return firstSymbols
}

// First computes the set of first symbols for a non-terminal node.
func (node *Node) First(visited map[*Node]bool) map[string]bool {
	var firstSymbols = map[string]bool{} // Set of all first terminal symbols

	if node == nil {
		return firstSymbols
	}

	// To avoid cycles, check if the node has already been visited
	if visited[node] {
		return firstSymbols
	}
	visited[node] = true

	if node.Terminal {
		// If it's a terminal node, add its symbol to the result
		firstSymbols[node.Tsym] = true
	} else if node.Nsym != nil && node.Nsym.Entry != nil {
		// If it's non-terminal, recurse on its Header Entry
		for sym := range node.Nsym.First(visited) {
			firstSymbols[sym] = true
		}
	}

	// Recurse on Alternative (Alt) nodes
	if node.Alt != nil {
		for sym := range node.Alt.First(visited) {
			firstSymbols[sym] = true
		}
	}

	return firstSymbols
}

var input = []byte(`A=x,(B).B=AC.C=[+A].`)

func getsym() {
	// get a symbol from the input stream
	if sym == '$' {
		return
	}
	read()
	write()
}

var sym rune

func read() {
	if len(input) == 0 {
		sym = '$'
		return
	}
	sym = rune(input[0])
	input = input[1:]
}

func write() {
	if sym == 0 || sym == '$' {
		return
	}
	fmt.Printf("%s", string(sym))
}

var symtab map[rune]*Header

// find locates a symbol in the symbol table.
// if it is not there, it is added to the symbol table
func find(sym rune) *Header {
	h, ok := symtab[sym]
	if !ok {
		h = &Header{
			Sym: string(sym),
		}
		symtab[sym] = h
	}
	return h
}

func werror() {
	fmt.Printf("	incorrect syntax\n")
	panic("syntax error")
}

func term(p, q, r *Node) {
	var a, b, c *Node

	var factor func(p, a *Node)
	factor = func(p, a *Node) {
		q = a
		for ('A' <= sym && sym <= 'Z') || sym == '[' || sym == 'ε' {
			factor(a.Suc, b)
			b.Alt = nil
			a = b
		}
		r = a
	}

}
