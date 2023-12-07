package kasane

import (
	"fmt"
	"strings"
	"testing"
)

func ExampleOverlayString() {
	base := ".......\n.......\n.......\n.......\n......."
	s := "xxx\nyyy\nzzz"

	out := OverlayString(base, s, 1, 3)
	fmt.Println(out)

	// Output:
	// .......
	// ...xxx.
	// ...yyy.
	// ...zzz.
	// .......
}

func TestOverlayString(t *testing.T) {
	tests := []struct {
		base, s   string
		top, left int
		want      string
	}{
		{
			base: join(
				"aa",
				"bb",
			),
			s:    "",
			top:  0,
			left: 0,
			want: join(
				"aa",
				"bb",
			),
		},
		{
			base: join(
				"aaaa",
				"bbbb",
				"cccc",
				"dddd",
			),
			s: join(
				"xx",
				"yy",
			),
			top:  1,
			left: 1,
			want: join(
				"aaaa",
				"bxxb",
				"cyyc",
				"dddd",
			),
		},
		{
			base: join(
				"aaaa",
				"bbbb",
				"cccc",
				"dddd",
			),
			s: join(
				"xxx",
				"yyy",
				"zzz",
			),
			top:  2,
			left: 2,
			want: join(
				"aaaa",
				"bbbb",
				"ccxx",
				"ddyy",
			),
		},
		{
			base: join(
				"aa",
				"bb",
			),
			s: join(
				"xx",
				"yy",
			),
			top:  -1,
			left: -1,
			want: join(
				"ya",
				"bb",
			),
		},
		{
			base: join(
				"aa",
				"bb",
			),
			s: join(
				"xx",
				"yy",
			),
			top:  3,
			left: 3,
			want: join(
				"aa",
				"bb",
			),
		},
		{
			base: join(
				"aa",
				"bbb",
				"cccc",
				"ddddd",
			),
			s: join(
				"xx",
				"yy",
				"z",
			),
			top:  1,
			left: 3,
			want: join(
				"aa",
				"bbb",
				"cccy",
				"dddzd",
			),
		},
		{
			base: join(
				"...",
				"....",
				".....",
				".....",
				"....",
			),
			s: join(
				"xxx",
				"y",
				"zzz",
			),
			top:  1,
			left: 1,
			want: join(
				"...",
				".xxx",
				".y...",
				".zzz.",
				"....",
			),
		},
		{
			base: join(
				"\x1b[31m.....\x1b[0m",
				"\x1b[31m.....\x1b[0m",
				"\x1b[31m.....\x1b[0m",
			),
			s: join(
				"xx",
				"yy",
			),
			top:  1,
			left: 1,
			want: join(
				"\x1b[31m.....\x1b[0m",
				"\x1b[31m.\x1b[0mxx\x1b[31m..\x1b[0m",
				"\x1b[31m.\x1b[0myy\x1b[31m..\x1b[0m",
			),
		},
		{
			base: join(
				".....",
				".....",
				".....",
			),
			s: join(
				"\x1b[41mxx\x1b[0m",
				"y\x1b[42my\x1b[0m",
			),
			top:  1,
			left: 2,
			want: join(
				".....",
				"..\x1b[41mxx\x1b[0m.",
				"..y\x1b[42my\x1b[0m.",
			),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			testOverlayString(t, test.base, test.s, test.top, test.left, test.want)
		})
	}
}

func TestOverlayString_WithPadding(t *testing.T) {
	tests := []struct {
		base, s   string
		top, left int
		pad       int
		want      string
	}{
		{
			base: join(
				"aa",
				"bbb",
				"cccc",
				"ddddd",
			),
			s: join(
				"xx",
				"yy",
				"z",
			),
			top:  1,
			left: 3,
			pad:  5,
			want: join(
				"aa   ",
				"bbbxx",
				"cccyy",
				"dddzd",
			),
		},
		{
			base: join(
				"...",
				"....",
				".....",
				".....",
				"....",
			),
			s: join(
				"xxx",
				"y",
				"zzz",
			),
			top:  1,
			left: 1,
			pad:  4,
			want: join(
				"... ",
				".xxx",
				".y...",
				".zzz.",
				"....",
			),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			testOverlayString(t, test.base, test.s, test.top, test.left, test.want, WithPadding(test.pad))
		})
	}
}

func join(ss ...string) string {
	return strings.Join(ss, "\n")
}

func testOverlayString(t *testing.T, base, s string, top, left int, want string, opts ...Option) {
	got := OverlayString(base, s, top, left, opts...)
	t.Logf("base string:\n%s", base)
	t.Logf("s string:\n%s", s)
	t.Logf("got string:\n%s", got)
	t.Logf("want string:\n%s", want)
	if got != want {
		t.Errorf("\ngot:\n%v\nwant:\n%v", got, want)
	}
}

func TestOverlaySingleLineString(t *testing.T) {
	tests := []struct {
		base, s string
		left    int
		want    string
	}{
		{
			base: "abc",
			s:    "",
			left: 0,
			want: "abc",
		},
		{
			base: "",
			s:    "xyz",
			left: 0,
			want: "",
		},
		{
			base: "abc",
			s:    "xyz",
			left: 0,
			want: "xyz",
		},
		{
			base: "abcde",
			s:    "xyz",
			left: 1,
			want: "axyze",
		},
		{
			base: "a",
			s:    "xyz",
			left: -1,
			want: "y",
		},
		{
			base: "abcd",
			s:    "xy",
			left: 3,
			want: "abcx",
		},
		{
			base: "ab",
			s:    "xy",
			left: 2,
			want: "ab",
		},
		{
			base: "abc",
			s:    "xyz",
			left: -2,
			want: "zbc",
		},
		{
			base: "abc",
			s:    "xy",
			left: -3,
			want: "abc",
		},
		{
			base: "abc",
			s:    "あ",
			left: 0,
			want: "あc",
		},
		{
			base: "abcdef",
			s:    "ああ",
			left: 1,
			want: "aああf",
		},
		{
			base: "abcdef",
			s:    "ああ",
			left: 3,
			want: "abcあ ",
		},
		{
			base: "abcdef",
			s:    "ああ",
			left: -1,
			want: " あdef",
		},
		{
			base: "abcdef",
			s:    "ああ",
			left: -3,
			want: " bcdef",
		},
		{
			base: "あいうえお",
			s:    "abcd",
			left: 4,
			want: "あいabcdお",
		},
		{
			base: "あいうえお",
			s:    "abcd",
			left: 3,
			want: "あ abcd お",
		},
		{
			base: "あいう",
			s:    "abcd",
			left: -1,
			want: "bcd う",
		},
		{
			base: "あいう",
			s:    "abcd",
			left: 3,
			want: "あ abc",
		},
		{
			base: "abcde",
			s:    "\x1b[31mxyz\x1b[0m",
			left: 1,
			want: "a\x1b[31mxyz\x1b[0me",
		},
		{
			base: "\x1b[41mabcde\x1b[0m",
			s:    "xyz",
			left: 1,
			want: "\x1b[41ma\x1b[0mxyz\x1b[41me\x1b[0m",
		},
		{
			base: "\x1b[41mabcde\x1b[0m",
			s:    "xyz",
			left: 3,
			want: "\x1b[41mabc\x1b[0mxy",
		},
		{
			base: "\x1b[41mabc\x1b[0m",
			s:    "xyz",
			left: 0,
			want: "xyz",
		},
		{
			base: "\x1b[41mabc\x1b[0m",
			s:    "",
			left: 0,
			want: "\x1b[41mabc\x1b[0m",
		},
		{
			base: "\x1b[41mabcde\x1b[0m",
			s:    "xyz",
			left: -1,
			want: "yz\x1b[41mcde\x1b[0m",
		},
		{
			base: "\x1b[1mabcde\x1b[0m",
			s:    "\x1b[4mxyz\x1b[0m",
			left: 1,
			want: "\x1b[1ma\x1b[0m\x1b[4mxyz\x1b[0m\x1b[1me\x1b[0m",
		},
		{
			base: "\x1b[1m\x1b[3m\x1b[42mabc\x1b[0mde",
			s:    "xy",
			left: 2,
			want: "\x1b[1m\x1b[3m\x1b[42mab\x1b[0mxye",
		},
		{
			base: "ab\x1b[1mcd\x1b[3mef\x1b[32mgh\x1b[0mij",
			s:    "\x1b[35mx\x1b[0m",
			left: 5,
			want: "ab\x1b[1mcd\x1b[3me\x1b[0m\x1b[35mx\x1b[0m\x1b[1m\x1b[3m\x1b[32mgh\x1b[0mij",
		},
		{
			base: "あ\x1b[43mいうえ\x1b[0mお",
			s:    "ab",
			left: 3,
			want: "あ ab \x1b[43mえ\x1b[0mお",
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			testOverlaySingleLineString(t, test.base, test.s, test.left, test.want)
		})
	}
}

func testOverlaySingleLineString(t *testing.T, base, s string, left int, want string) {
	k := new()
	got := k.overlaySingleLineString(base, s, left)
	t.Logf("base string: `%s`", base)
	t.Logf("   s string: `%s`", s)
	t.Logf("       left: %d", left)
	t.Logf("want string: `%s`", want)
	t.Logf(" got string: `%s`", got)
	if got != want {
		t.Errorf("\ngot:  `%v`\nwant: `%v`", got, want)
	}
}

func TestStringWidth(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{s: "", want: 0},
		{s: "abc", want: 3},
		{s: "あいう", want: 6},
		{s: "aあbい", want: 6},
		{s: "\x1b[31m\x1b[0m", want: 0},
		{s: "\033[31mabc\033[0m", want: 3},
		{s: "\x1b[38;5;123mabc\x1b[0mあい\x1b[1mう\x1b[0m", want: 9},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			k := new()
			got := k.stringWidth(test.s)
			if got != test.want {
				t.Errorf("s=%v, got=%v, want=%v", test.s, got, test.want)
			}
		})
	}
}

func TestToCells(t *testing.T) {
	tests := []struct {
		s    string
		want []*cell
	}{
		{s: "", want: []*cell{}},
		{s: "abc", want: []*cell{
			{r: 'a', double: false, head: true, csis: []string{}},
			{r: 'b', double: false, head: true, csis: []string{}},
			{r: 'c', double: false, head: true, csis: []string{}},
		}},
		{s: "あいう", want: []*cell{
			{r: 'あ', double: true, head: true, csis: []string{}},
			{r: 'あ', double: true, head: false, csis: []string{}},
			{r: 'い', double: true, head: true, csis: []string{}},
			{r: 'い', double: true, head: false, csis: []string{}},
			{r: 'う', double: true, head: true, csis: []string{}},
			{r: 'う', double: true, head: false, csis: []string{}},
		}},
		{s: "aあbい", want: []*cell{
			{r: 'a', double: false, head: true, csis: []string{}},
			{r: 'あ', double: true, head: true, csis: []string{}},
			{r: 'あ', double: true, head: false, csis: []string{}},
			{r: 'b', double: false, head: true, csis: []string{}},
			{r: 'い', double: true, head: true, csis: []string{}},
			{r: 'い', double: true, head: false, csis: []string{}},
		}},
		{s: "\x1b[31mabc\x1b[0m", want: []*cell{
			{r: 'a', double: false, head: true, csis: []string{"\x1b[31m"}},
			{r: 'b', double: false, head: true, csis: []string{"\x1b[31m"}},
			{r: 'c', double: false, head: true, csis: []string{"\x1b[31m"}},
		}},
		{s: "a\x1b[31mb\x1b[32mc\x1b[33md\x1b[0me", want: []*cell{
			{r: 'a', double: false, head: true, csis: []string{}},
			{r: 'b', double: false, head: true, csis: []string{"\x1b[31m"}},
			{r: 'c', double: false, head: true, csis: []string{"\x1b[31m", "\x1b[32m"}},
			{r: 'd', double: false, head: true, csis: []string{"\x1b[31m", "\x1b[32m", "\x1b[33m"}},
			{r: 'e', double: false, head: true, csis: []string{}},
		}},
		{s: "a\x1b[31mb\x1b[0mc\x1b[32md\x1b[0me", want: []*cell{
			{r: 'a', double: false, head: true, csis: []string{}},
			{r: 'b', double: false, head: true, csis: []string{"\x1b[31m"}},
			{r: 'c', double: false, head: true, csis: []string{}},
			{r: 'd', double: false, head: true, csis: []string{"\x1b[32m"}},
			{r: 'e', double: false, head: true, csis: []string{}},
		}},
		{s: "\x1b[38;5;123mabc\x1b[0mあい\x1b[1mう\x1b[0m", want: []*cell{
			{r: 'a', double: false, head: true, csis: []string{"\x1b[38;5;123m"}},
			{r: 'b', double: false, head: true, csis: []string{"\x1b[38;5;123m"}},
			{r: 'c', double: false, head: true, csis: []string{"\x1b[38;5;123m"}},
			{r: 'あ', double: true, head: true, csis: []string{}},
			{r: 'あ', double: true, head: false, csis: []string{}},
			{r: 'い', double: true, head: true, csis: []string{}},
			{r: 'い', double: true, head: false, csis: []string{}},
			{r: 'う', double: true, head: true, csis: []string{"\x1b[1m"}},
			{r: 'う', double: true, head: false, csis: []string{"\x1b[1m"}},
		}},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			k := new()
			got := k.toCells(test.s)
			if len(got) != len(test.want) {
				t.Errorf("s=%v, len(got)=%v, len(want)=%v", test.s, len(got), len(test.want))
			}
			for i := 0; i < len(got); i++ {
				if got[i].r != test.want[i].r {
					t.Errorf("s=%v, got[%v].r=%v, want[%v].r=%v", test.s, i, got[i].r, i, test.want[i].r)
				}
				if got[i].double != test.want[i].double {
					t.Errorf("s=%v, got[%v].double=%v, want[%v].double=%v", test.s, i, got[i].double, i, test.want[i].double)
				}
				if got[i].head != test.want[i].head {
					t.Errorf("s=%v, got[%v].head=%v, want[%v].head=%v", test.s, i, got[i].head, i, test.want[i].head)
				}
				if len(got[i].csis) != len(test.want[i].csis) {
					t.Errorf("s=%v, len(got[%v].csis)=%v, len(want[%v].csis)=%v", test.s, i, len(got[i].csis), i, len(test.want[i].csis))
				} else {
					for j := range got[i].csis {
						if got[i].csis[j] != test.want[i].csis[j] {
							t.Errorf("s=%v, got[%v].csis[%v]=%v, want[%v].csis[%v]=%v", test.s, i, j, got[i].csis[j], i, j, test.want[i].csis[j])
						}
					}
				}
			}
		})
	}
}
