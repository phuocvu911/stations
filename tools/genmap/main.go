package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

func main() {
	out := "maps/big.map"
	n, m := 10000, 30000

	args := os.Args[1:]
	if len(args) > 0 {
		out = args[0]
	}
	if len(args) > 1 {
		n = mustAtoi(args[1])
	}
	if len(args) > 2 {
		m = mustAtoi(args[2])
	}

	// Max unique undirected edges without self-loops.
	if maxEdges := n * (n - 1) / 2; m > maxEdges {
		fmt.Fprintf(os.Stderr, "Error: %d connections is more than %d unique pairs for %d stations\n", m, maxEdges, n)
		os.Exit(1)
	}

	f, err := os.Create(out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	rng := rand.New(rand.NewSource(42))

	// Stations: each on a unique coordinate in a 100-wide grid.
	fmt.Fprintln(w, "stations:")
	for i := 0; i < n; i++ {
		fmt.Fprintf(w, "station%d,%d,%d\n", i, i%100, i/100)
	}

	// Connections: unique undirected pairs.
	fmt.Fprintln(w)
	fmt.Fprintln(w, "connections:")
	seen := make(map[[2]int]struct{}, m)
	for count := 0; count < m; {
		a := rng.Intn(n)
		b := rng.Intn(n)
		if a == b {
			continue
		}
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		fmt.Fprintf(w, "station%d-station%d\n", a, b)
		count++
	}

	fmt.Fprintf(os.Stderr, "wrote %s: %d stations, %d connections\n", out, n, m)
}

func mustAtoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %q is not an integer\n", s)
		os.Exit(1)
	}
	return v
}
