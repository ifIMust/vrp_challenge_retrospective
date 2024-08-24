package input

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/ifIMust/vrp_challenge/common"
)

// ReadFile parses a VRP problem from a file, and outputs a
// slice of Loads if successful.
func ReadFile(fileName string) ([]*common.Load, error) {
	loads := make([]*common.Load, 0)
	file, err := os.Open(fileName)
	if err != nil {
		return loads, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Process line. Errors reading a line are silently ignored,
		// which discards the expected header row.
		load, err := loadFromLine(scanner.Text())
		if err == nil {
			loads = append(loads, load)
		}
	}
	err = scanner.Err()
	return loads, err
}

// loadFromLine creates a Load from a line of formatted text.
// Expected input style is space-delimited with 3 fields, e.g.:
// "1 (12.34,56.78) (12.34,56.78)"
func loadFromLine(line string) (*common.Load, error) {
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
	return common.NewLoad(index, pickup, dropoff), nil
}

// parseCoords creates a Location from formatted text.
// Expected input style: "(12.34,56.78)"
func parseCoords(coords string) (*common.Location, error) {
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
	return common.NewLocation(x, y), nil
}
