package input

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

type Location struct {
	X float64
	Y float64
}

type Load struct {
	Index   int
	Pickup  *Location
	Dropoff *Location
}

func ReadFile(fileName string) []*Load {
	loads := make([]*Load, 0)
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
		return loads
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Process line
		load, err := loadFromLine(scanner.Text())
		if err == nil {
			loads = append(loads, load)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return loads
}

func loadFromLine(line string) (*Load, error) {
	fields := strings.Split(line, " ")
	if len(fields) != 3 {
		return nil, errors.New("unexpected number of fields")
	}

	index, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, errors.New("index was not an integer")
	}

	pickup, err := parseCoords(fields[1])
	if err != nil {
		return nil, errors.New("invalid pickup coordinates")
	}

	dropoff, err := parseCoords(fields[2])
	if err != nil {
		return nil, errors.New("invalid dropoff coordinates")
	}

	return &Load{Index: index, Pickup: pickup, Dropoff: dropoff}, nil
}

// Expected input style: "(12.34,56.78)"
func parseCoords(coords string) (*Location, error) {
	coords = strings.ReplaceAll(coords, "(", "")
	coords = strings.ReplaceAll(coords, ")", "")
	fields := strings.Split(coords, ",")
	if len(fields) != 2 {
		return nil, errors.New("unexpected coordinate format")
	}

	x, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return nil, errors.New("X coordinate could not be interpreted as a float.")
	}

	y, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, errors.New("Y coordinate could not be interpreted as a float.")
	}
	return &Location{X: x, Y: y}, nil
}
