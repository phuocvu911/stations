package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type simTrain struct {
	path int
	pos  int
	done bool
}

// simulate runs the turns. Each route is a pipeline: one new train enters it
// per turn while every train already on it advances one station, so trains
// stay one block apart and no track is reused within a turn. visit is called
// once per turn with that turn's movements and every train created so far
// (train Tn is trains[n-1]).
func simulate(paths [][]int, counts []int, names []string, visit func(movements []string, trains []simTrain)) {
	var trains []simTrain
	dispatched := make([]int, len(paths))
	for {
		var line []string
		for ti := range trains {
			t := &trains[ti]
			if t.done {
				continue
			}
			t.pos++
			line = append(line, fmt.Sprintf("T%d-%s", ti+1, names[paths[t.path][t.pos]]))
			if t.pos == len(paths[t.path])-1 {
				t.done = true
			}
		}
		for pi := range paths {
			if dispatched[pi] < counts[pi] {
				dispatched[pi]++
				trains = append(trains, simTrain{path: pi, pos: 1, done: len(paths[pi]) == 2})
				line = append(line, fmt.Sprintf("T%d-%s", len(trains), names[paths[pi][1]]))
			}
		}
		if len(line) == 0 {
			break
		}
		visit(line, trains)
	}
}

func printSchedule(net *Network, paths [][]int, counts []int) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	simulate(paths, counts, net.names, func(movements []string, _ []simTrain) {
		fmt.Fprintln(w, strings.Join(movements, " "))
	})
}
