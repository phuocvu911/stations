package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	// start := time.Now()
	// defer func() {
	// 	fmt.Printf("Execution time: %s\n", time.Since(start))
	// }()

	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
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
		return fmt.Errorf("Error: number of trains (%d) is not a valid positive integer", numTrains)
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

	bestPaths, trainsEachPath, ok := planMovements(net, startStation, endStation, numTrains)
	if !ok {
		return fmt.Errorf("Error: no path exists between %q and %q", startName, endName)
	}
	//fmt.Println(bestPaths, trainsEachPath)

	printSchedule(net, bestPaths, trainsEachPath)
	return nil
}
