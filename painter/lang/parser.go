package lang

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

var (
	ErrEmptyLine          = errors.New("empty line")
	ErrUnknownCommand     = errors.New("unknown command")
	ErrInsufficientParams = errors.New("insufficient number of parameters")
)

type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)
	var res []painter.Operation
	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := parse(commandLine)
		if err != nil {
			return nil, err
		}

		res = append(res, op)
	}
	return res, nil
}

func parse(line string) (painter.Operation, error) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, ErrEmptyLine
	}

	cmd, params := parts[0], parts[1:]
	switch cmd {
	case "white":
		return painter.WhiteFill, nil
	case "green":
		return painter.GreenFill, nil
	case "update":
		return painter.Update, nil
	case "bgrect":
		if len(params) < 4 {
			return nil, fmt.Errorf("%w: have: %d, want: %d", ErrInsufficientParams, len(params), 4)
		}
		p1, err := parsePoint(params[0], params[1])
		if err != nil {
			return nil, err
		}
		p2, err := parsePoint(params[2], params[3])
		if err != nil {
			return nil, err
		}
		return painter.BgRect(painter.Rectangle{Min: p1, Max: p2}), nil
	case "figure":
		if len(params) < 2 {
			return nil, fmt.Errorf("%w: have: %d, want: %d", ErrInsufficientParams, len(params), 2)
		}
		p, err := parsePoint(params[0], params[1])
		if err != nil {
			return nil, err
		}
		return painter.Figure(p), nil
	case "move":
		if len(params) < 2 {
			return nil, fmt.Errorf("%w: have: %d, want: %d", ErrInsufficientParams, len(params), 2)
		}
		p, err := parsePoint(params[0], params[1])
		if err != nil {
			return nil, err
		}
		return painter.Move(p), nil
	case "reset":
		return painter.Reset, nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnknownCommand, cmd)
	}
}

func parsePoint(x, y string) (painter.Point, error) {
	xf, err := strconv.ParseFloat(x, 32)
	if err != nil {
		return painter.Point{}, err
	}
	yf, err := strconv.ParseFloat(y, 32)
	if err != nil {
		return painter.Point{}, err
	}
	return painter.Pt(float32(xf), float32(yf)), nil
}
