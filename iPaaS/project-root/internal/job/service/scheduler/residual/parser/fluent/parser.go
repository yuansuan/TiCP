package fluent

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual/parser"
)

var fluentIterValuesRegx = regexp.MustCompile(`^\d+\s+.*`)

var iterWords = []string{"iter", "continuity", "x-velocity"}

// Record defines the record
type Record struct {
	Series map[string]float64
	Iter   int
}

// Context defines the fluent context
type Context struct {
	lastRecord    *Record
	iterationLine string
	varNames      []string
	iter          int
}

// Parser parses the result from a residual log
type Parser struct {
}

// Parse parses the result from a residual log
func (p *Parser) Parse(r io.Reader) (res *parser.Result, err error) {
	// prevent system's panic
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
			res = nil
		}
	}()

	scanner := bufio.NewScanner(r)

	ctx := &Context{
		iterationLine: "",
		iter:          -1,
	}

	res = &parser.Result{
		Series: map[string][]float64{},
		XVar:   []string{"iter"},
	}

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case p.isIterLine(ctx, line):
			ctx.iterationLine = line
			ctx.varNames = p.extractVarNames(ctx, line)
		case ctx.iterationLine != "" && p.isIterValue(ctx, line):
			record, err := p.readIterValue(ctx, line)
			if err != nil {
				continue
			}
			ctx.lastRecord = record
			ctx.iter = record.Iter
			p.mergeResult(ctx, res, record)
		}
	}

	return res, scanner.Err()
}

func (p *Parser) isIterLine(ctx *Context, line string) bool {
	if !strings.HasPrefix(strings.TrimSpace(line), "iter") {
		return false
	}

	fields := strings.Fields(line)

	contain := func(arr []string, str string) bool {
		for _, v := range arr {
			if v == str {
				return true
			}
		}
		return false
	}

	for _, v := range iterWords {
		if !contain(fields, v) {
			return false
		}
	}
	return true
}

func (p *Parser) isIterValue(ctx *Context, line string) bool {
	return fluentIterValuesRegx.MatchString(strings.TrimSpace(line))
}

func (p *Parser) readIterValue(ctx *Context, line string) (*Record, error) {
	values, err := p.extractVarValues(ctx, line)
	if err != nil {
		return nil, err
	}

	record := &Record{
		Series: map[string]float64{},
		Iter:   int(values[0]),
	}

	for i, v := range values {
		record.Series[ctx.varNames[i]] = v
	}

	return record, nil
}

func (p *Parser) extractVarValues(ctx *Context, line string) ([]float64, error) {
	values := strings.Fields(line)
	values = values[:len(values)-2]

	var floatVals []float64

	for _, v := range values {
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		floatVals = append(floatVals, fv)
	}
	return floatVals, nil
}

func (p *Parser) extractVarNames(ctx *Context, line string) []string {
	names := strings.Fields(line)
	return names[:len(names)-1]
}

func (p *Parser) mergeResult(ctx *Context, res *parser.Result, record *Record) {
	for _, name := range ctx.varNames {
		res.Series[name] = append(res.Series[name], record.Series[name])
	}
}
