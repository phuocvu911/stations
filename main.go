package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	args, _ = extractBonusFlag(args)
	if len(args) != 4 {
		return fmt.Errorf("Error: incorrect number of command line arguments\nUsage: go run . [path to network map] [start station] [end station] [number of trains]")
	}
	mapPath, startName, endName := args[0], args[1], args[2]
	numTrains, err := strconv.Atoi(args[3])
	if err != nil || numTrains <= 0 {
		return fmt.Errorf("Error: number of trains is not a valid positive integer")
	}
	net, err := parseNetwork(mapPath)
	if err != nil {
		return err
	}
	//fmt.Printf("%+v\n", net)
	startStation, ok := net.index[startName]
	if !ok {
		return fmt.Errorf("Error: start station %q does not exist", startName)
	}
	endStation, ok := net.index[endName]
	if !ok {
		return fmt.Errorf("Error: end station %q does not exist", endName)
	}
	if startStation == endStation {
		return fmt.Errorf("Error: start and end station are the same")
	}
	paths, counts, ok := planMovements(net, startStation, endStation, numTrains)
	if !ok {
		return fmt.Errorf("Error: no path exists between %q and %q", startName, endName)
	}

	printSchedule(net, paths, counts)
	return nil
}
