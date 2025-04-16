package lang

import (
	"errors"
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

func TestParse(t *testing.T) {
	type want struct {
		ops []painter.Operation
		err error
	}
	cases := []struct {
		name  string
		input string
		want  want
	}{
		{
			name:  "one operation",
			input: "white",
			want: want{
				ops: []painter.Operation{
					painter.WhiteFill,
				},
			},
		},
		{
			name:  "all operations",
			input: "white \n green \n update \n bgrect 0.5 0.75 0.95 1 \n figure 0.01 0.01 \n move 0.69 0.69 \n reset",
			want: want{
				ops: []painter.Operation{
					painter.WhiteFill,
					painter.GreenFill,
					painter.Update,
					painter.BgRect(painter.Rect(0.5, 0.75, 0.95, 1)),
					painter.Figure(painter.Pt(0.01, 0.01)),
					painter.Move(painter.Pt(0.69, 0.69)),
					painter.Reset,
				},
			},
		},
		{
			name:  "empty line",
			input: "white \n  \n green",
			want: want{
				err: ErrEmptyLine,
			},
		},
		{
			name:  "insufficient number of parameters",
			input: "bgrect 0.5 0.75",
			want: want{
				err: ErrInsufficientParams,
			},
		},
		{
			name:  "bad command",
			input: "badCommand",
			want: want{
				err: ErrUnknownCommand,
			},
		},
	}

	var p Parser

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ops, err := p.Parse(strings.NewReader(tc.input))
			if err != nil {
				if !errors.Is(err, tc.want.err) {
					t.Fatalf("Parse(%s) = err: %v, want: %v", tc.input, err, tc.want.err)
				}
			}

			if len(ops) != len(tc.want.ops) {
				t.Fatalf("Parse(%s) len(ops): %d, want: %d", tc.input, len(ops), len(tc.want.ops))
			}
			for i, op := range ops {
				wantOp := tc.want.ops[i]
				if !isOpsEqual(op, wantOp) {
					t1 := painter.MockState()
					r1 := op.Do(&t1)
					t2 := painter.MockState()
					r2 := wantOp.Do(&t2)
					t.Fatalf("Parse(\n%s\n) failed on %d operation: ready and texture state after ops: (%t, %v), want: (%t, %v)", tc.input, i, r1, t1, r2, t2)
				}
			}
		})
	}
}

func isOpsEqual(op1, op2 painter.Operation) bool {
	s1 := painter.MockState()
	r1 := op1.Do(&s1)
	s2 := painter.MockState()
	r2 := op2.Do(&s2)
	return r1 == r2 && s1.Equal(s2)
}
