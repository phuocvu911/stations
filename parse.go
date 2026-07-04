package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const maxStations = 10000

// Network is the parsed rail network. Stations are referred to by their
// index into stationNames; conns holds each undirected track exactly once.
type Network struct {
	stationNames []string
	//xs, ys []int
	index       map[string]int //station name -> index lookup
	connections [][2]int       //connections between stations in the form of index
	adj         [][]int        //adjacent stations in adj[v], for example adj[0]= {1,3,4} mean station 0 is next to AND have connection with station 1,3,4
}

type rawConn struct {
	a, b string
	line int
}

func parseNetwork(path string) (*Network, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error: cannot read network map file %q", path)
	}
	net := &Network{index: make(map[string]int)}
	coords := make(map[[2]int]string)

	var rawConns []rawConn
	section := ""
	sawStations, sawConns := false, false
	for i, raw := range strings.Split(string(data), "\n") {
		lineNo := i + 1
		line := raw
		//skip everything after '#'
		if j := strings.IndexByte(line, '#'); j >= 0 {
			line = line[:j]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "stations:" {
			if sawStations {
				return nil, fmt.Errorf("Error: duplicate \"stations:\" section on line %d", lineNo)
			}
			sawStations, section = true, "stations"
			continue
		}
		if line == "connections:" {
			if sawConns {
				return nil, fmt.Errorf("Error: duplicate \"connections:\" section on line %d", lineNo)
			}
			sawConns, section = true, "connections"
			continue
		}
		switch section {
		case "stations":
			parts := strings.Split(line, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("Error: invalid station definition on line %d: %q", lineNo, line)
			}
			name := strings.TrimSpace(parts[0])
			if !validStationName(name) {
				return nil, fmt.Errorf("Error: invalid station name on line %d: %q", lineNo, name)
			}
			x, okX := parsePositiveInt(strings.TrimSpace(parts[1]))
			y, okY := parsePositiveInt(strings.TrimSpace(parts[2]))
			if !okX || !okY {
				return nil, fmt.Errorf("Error: station %q has a coordinate which is not a valid positive integer (line %d)", name, lineNo)
			}
			//check duplicate station
			if _, dup := net.index[name]; dup {
				return nil, fmt.Errorf("Error: duplicate station name %q (line %d)", name, lineNo)
			}
			//check if 2 stations have the same coordinates
			if otherStation, taken := coords[[2]int{x, y}]; taken {
				return nil, fmt.Errorf("Error: stations %q and %q exist at the same coordinates %d,%d", otherStation, name, x, y)
			}
			coords[[2]int{x, y}] = name
			net.index[name] = len(net.stationNames) //name -> index lookup, 0,1,2,3,..., because it set before appending name
			net.stationNames = append(net.stationNames, name)
			// net.xs = append(net.xs, x)
			// net.ys = append(net.ys, y)
			if len(net.stationNames) > maxStations {
				return nil, fmt.Errorf("Error: the map contains more than %d stations", maxStations)
			}
		case "connections":
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Error: invalid connection declaration on line %d: %q", lineNo, line)
			}
			rawConns = append(rawConns, rawConn{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), lineNo})
		default:
			return nil, fmt.Errorf("Error: line %d does not belong to a \"stations:\" or \"connections:\" section: %q", lineNo, line)
		}
	}
	if !sawStations {
		return nil, fmt.Errorf("Error: the map does not contain a \"stations:\" section")
	}
	if !sawConns {
		return nil, fmt.Errorf("Error: the map does not contain a \"connections:\" section")
	}
	net.adj = make([][]int, len(net.stationNames))
	connSeen := make(map[[2]int]bool)
	for _, rc := range rawConns {
		ia, okA := net.index[rc.a]
		if !okA {
			return nil, fmt.Errorf("Error: connection on line %d refers to station %q which does not exist", rc.line, rc.a)
		}
		ib, okB := net.index[rc.b]
		if !okB {
			return nil, fmt.Errorf("Error: connection on line %d refers to station %q which does not exist", rc.line, rc.b)
		}
		if ia == ib {
			return nil, fmt.Errorf("Error: connection on line %d connects station %q to itself", rc.line, rc.a)
		}
		key := [2]int{min(ia, ib), max(ia, ib)}
		if connSeen[key] {
			return nil, fmt.Errorf("Error: duplicate connection between %q and %q (line %d)", rc.a, rc.b, rc.line)
		}
		connSeen[key] = true
		net.connections = append(net.connections, [2]int{ia, ib})
		net.adj[ia] = append(net.adj[ia], ib)
		net.adj[ib] = append(net.adj[ib], ia)
	}
	return net, nil
}

func validStationName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' {
			return false
		}
	}
	return true
}

func parsePositiveInt(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 || strings.HasPrefix(s, "+") {
		return 0, false
	}
	return n, true
}
