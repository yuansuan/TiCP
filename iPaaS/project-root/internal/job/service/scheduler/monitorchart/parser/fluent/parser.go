package fluent

import (
	"bufio"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/monitorchart/parser"
)

// Parser parses the monitor chart
type Parser struct{}

var monitorIterValuesRegx = regexp.MustCompile(`^\d+\s+.*`)

// Parse parses the monitor chart data
func (p *Parser) Parse(file string, r io.Reader) (resMap map[string]*parser.Result, err error) {
	// prevent system's panic
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			resMap = nil
		}
	}()

	key := strings.Split(filepath.Base(file), ".")[0]

	result := &parser.Result{
		Key:   key,
		Items: []*parser.ParseItem{},
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if monitorIterValuesRegx.MatchString(strings.TrimSpace(line)) {
			item := p.readParseItem(key, line)
			result.Items = append(result.Items, item)
		}
	}

	// Assuming resMap is initialized before calling Parse
	resMap = make(map[string]*parser.Result)
	resMap[key] = result

	return resMap, scanner.Err()
}

func (p *Parser) readParseItem(key, line string) *parser.ParseItem {
	var iteration, value float64
	var err error

	values := strings.Fields(line)

	if len(values) >= 2 {
		_, err = fmt.Sscanf(values[0], "%f", &iteration)
		if err != nil {
			panic(err)
		}

		_, err = fmt.Sscanf(values[1], "%f", &value)
		if err != nil {
			panic(err)
		}
	}

	return &parser.ParseItem{
		Key:       key,
		Iteration: float64(iteration),
		Value:     float64(value),
	}
}
