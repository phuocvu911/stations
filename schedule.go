package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type simTrain struct {
	path int  //index of the route in paths slice
	pos  int  //index of the station in the route,refered as paths[path]
	done bool //true if the train has reached end station
}

// simulate runs the turns. Each route is a pipeline: one new train enters it
// per turn while every train already on it advances one station, so trains
// stay one block apart and no track is reused within a turn. visit is called
// once per turn with that turn's movements and every train created so far
// (train Tn is trains[n-1]).
func simulate(bestPaths [][]int, trainsEachPath []int, names []string) (lines [][]string) {
	var trains []simTrain

	//dispatched[i] tracks how many trains have already been sent onto paths[i]
	dispatched := make([]int, len(bestPaths))

	for {
		var line []string

		//advance existing trains
		for i := range trains {
			train := &trains[i] //get a pointer to the train so we can modify it
			if train.done {
				continue
			}

			train.pos++ //move the train forward one station along its route

			line = append(line, fmt.Sprintf("T%d-%s", i+1, names[bestPaths[train.path][train.pos]]))

			// If the new position is the last index of the path, the train has arrived
			if train.pos == len(bestPaths[train.path])-1 {
				train.done = true
			}
		}

		//launch new trains
		for i := range bestPaths {
			if dispatched[i] < trainsEachPath[i] {
				dispatched[i]++
				trains = append(trains, simTrain{path: i, pos: 1, done: len(bestPaths[i]) == 2})
				line = append(line, fmt.Sprintf("T%d-%s", len(trains), names[bestPaths[i][1]]))
			}
		}

		//if no trains moved or were launched this turn, we're done
		if len(line) == 0 {
			break
		}
		lines = append(lines, line)
	}
	return lines
}

func printSchedule(net *Network, paths [][]int, counts []int) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	lines := simulate(paths, counts, net.stationNames)
	for _, line := range lines {
		fmt.Fprintln(w, strings.Join(line, " "))
	}
}
