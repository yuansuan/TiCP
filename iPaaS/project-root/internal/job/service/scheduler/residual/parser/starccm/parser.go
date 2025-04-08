package starccm

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/scheduler/residual/parser"
)

var (
	starCCMIterValuesRegx = regexp.MustCompile(`^\d+\s+.*`)

	// matchs: TimeStep  596: Time  5.960000e-01
	starCCMTimeRegexp = regexp.MustCompile(`TimeStep\s+(\d+): Time\s+(\S+)`)
)

const (
	starCCMTimeStepVar = "TimeStep_Time"
	starCCMIterVar     = "Iteration"
)

// Context defines the starccm context
type Context struct {
	haveTimeStep  bool
	lastTimeStep  float64
	lastRecord    *Record
	iterationLine string
	varNames      []string
	iter          int
	records       []*Record
	stepTime      bool
	count         int
}

// Record defines the record
type Record struct {
	Series map[string]float64
	Iter   int
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

	// 设置缓冲区为1MB
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	ctx := &Context{
		haveTimeStep:  false,
		lastTimeStep:  0,
		lastRecord:    nil,
		iterationLine: "",
		iter:          -1,
		stepTime:      false,
		count:         0,
	}

	res = &parser.Result{
		Series: map[string][]float64{},
		XVar:   []string{},
	}

	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case p.isIterLine(ctx, line):
			ctx.iterationLine = line
		case p.isTimeStepLine(ctx, line):
			ctx.haveTimeStep = true
			ctx.lastTimeStep = p.readTimeStep(ctx, line)
			if ctx.count != 0 {
				ctx.records = ctx.records[0 : len(ctx.records)-ctx.count]
				ctx.count = 0
			}
			ctx.stepTime = true
		case ctx.iterationLine != "":
			if ctx.stepTime {
				record, err := p.readIterValue(ctx, line)
				if err != nil {
					continue
				}
				ctx.lastRecord = record
				ctx.records = append(ctx.records, record)
				ctx.stepTime = false
				continue
			}
			if p.isIterValue(ctx, line) {
				record, err := p.readIterValue(ctx, line)
				if err != nil {
					continue
				}
				ctx.lastRecord = record
				ctx.records = append(ctx.records, record)
				ctx.iter = record.Iter
				ctx.count++
			}
		}
	}

	p.mergeResult(ctx, res, ctx.records)

	res.XVar = []string{}
	if ctx.haveTimeStep {
		res.XVar = append(res.XVar, starCCMTimeStepVar)
	} else {
		res.XVar = append(res.XVar, starCCMIterVar)
	}

	return res, scanner.Err()
}

func (p *Parser) mergeResult(ctx *Context, res *parser.Result, records []*Record) {
	for _, record := range records {
		for _, name := range ctx.varNames {
			res.Series[name] = append(res.Series[name], record.Series[name])
		}

		if ctx.haveTimeStep {
			res.Series[starCCMTimeStepVar] = append(res.Series[starCCMTimeStepVar], record.Series[starCCMTimeStepVar])
		}
	}
}

func (p *Parser) isIterLine(ctx *Context, line string) bool {
	// start with "Iteration"
	ret := strings.HasPrefix(strings.TrimSpace(line), starCCMIterVar)
	if ret == false {
		return false
	}

	// and Iteration line contain "Continuity"
	arr := strings.Fields(line)
	for _, field := range arr {
		if field == "Continuity" || field == "Energy" {
			return true
		}
	}
	return false
}

func (p *Parser) isTimeStepLine(ctx *Context, line string) bool {
	// match: TimeStep  596: Time  5.960000e-01
	return starCCMTimeRegexp.MatchString(strings.TrimRightFunc(line, unicode.IsSpace))
}

func (p *Parser) readTimeStep(ctx *Context, line string) float64 {
	match := starCCMTimeRegexp.FindStringSubmatch(line)
	if len(match) < 3 {
		panic("submatch should len(submatch)==3")
	}

	// ["TimeStep  593: Time  5.930000e-01", "593", "5.930000e-01"]
	time, err := strconv.ParseFloat(match[2], 64)
	if err != nil {
		panic(err)
	}

	return time
}

func (p *Parser) isIterValue(ctx *Context, line string) bool {
	return starCCMIterValuesRegx.MatchString(strings.TrimSpace(line))
}

func (p *Parser) readIterValue(ctx *Context, line string) (*Record, error) {
	values, err := p.extractVarValues(ctx, line)
	if err != nil {
		return nil, err
	}
	if len(values) > len(ctx.varNames) {
		ctx.varNames = p.extractVarName(ctx, line)
	}

	record := &Record{
		Series: map[string]float64{},
		Iter:   int(values[0]),
	}

	for i, v := range values {
		record.Series[ctx.varNames[i]] = v
	}

	if ctx.haveTimeStep {
		record.Series[starCCMTimeStepVar] = float64(ctx.lastTimeStep)
	}

	return record, nil
}

func (p *Parser) extractVarName(ctx *Context, line string) []string {
	valueRegx := regexp.MustCompile(`\s*\S+`)
	startEnds := valueRegx.FindAllStringIndex(line, -1)

	symbols := []string{}
	// we skip the first symbol `Iteration`
	for _, se := range startEnds {
		symbols = append(symbols, strings.TrimSpace(ctx.iterationLine[se[0]:se[1]]))
	}

	return symbols
}

func (p *Parser) extractVarValues(ctx *Context, line string) ([]float64, error) {
	var floatVals []float64
	values := strings.Fields(line)
	for _, v := range values {
		fv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		floatVals = append(floatVals, fv)
	}
	return floatVals, nil
}
