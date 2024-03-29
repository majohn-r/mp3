package files_test

import (
	"mp3/internal/files"
	"testing"
)

func Test_isIllegalRuneForFileNames(t *testing.T) {
	const fnName = "isIllegalRuneForFileNames()"
	type args struct {
		r rune
	}
	tests := map[string]struct {
		args
		want bool
	}{
		"_0":    {args: args{r: 0}, want: true},
		"_1":    {args: args{r: 1}, want: true},
		"_2":    {args: args{r: 2}, want: true},
		"_3":    {args: args{r: 3}, want: true},
		"_4":    {args: args{r: 4}, want: true},
		"_5":    {args: args{r: 5}, want: true},
		"_6":    {args: args{r: 6}, want: true},
		"_7":    {args: args{r: 7}, want: true},
		"_8":    {args: args{r: 8}, want: true},
		"_9":    {args: args{r: 9}, want: true},
		"10":    {args: args{r: 10}, want: true},
		"11":    {args: args{r: 11}, want: true},
		"12":    {args: args{r: 12}, want: true},
		"13":    {args: args{r: 13}, want: true},
		"14":    {args: args{r: 14}, want: true},
		"15":    {args: args{r: 15}, want: true},
		"16":    {args: args{r: 16}, want: true},
		"17":    {args: args{r: 17}, want: true},
		"18":    {args: args{r: 18}, want: true},
		"19":    {args: args{r: 19}, want: true},
		"20":    {args: args{r: 20}, want: true},
		"21":    {args: args{r: 21}, want: true},
		"22":    {args: args{r: 22}, want: true},
		"23":    {args: args{r: 23}, want: true},
		"24":    {args: args{r: 24}, want: true},
		"25":    {args: args{r: 25}, want: true},
		"26":    {args: args{r: 26}, want: true},
		"27":    {args: args{r: 27}, want: true},
		"28":    {args: args{r: 28}, want: true},
		"29":    {args: args{r: 29}, want: true},
		"30":    {args: args{r: 30}, want: true},
		"31":    {args: args{r: 31}, want: true},
		"<":     {args: args{r: '<'}, want: true},
		">":     {args: args{r: '>'}, want: true},
		":":     {args: args{r: ':'}, want: true},
		"\"":    {args: args{r: '"'}, want: true},
		"/":     {args: args{r: '/'}, want: true},
		"\\":    {args: args{r: '\\'}, want: true},
		"|":     {args: args{r: '|'}, want: true},
		"?":     {args: args{r: '?'}, want: true},
		"*":     {args: args{r: '*'}, want: true},
		"!":     {args: args{r: '!'}, want: false},
		"#":     {args: args{r: '#'}, want: false},
		"$":     {args: args{r: '$'}, want: false},
		"&":     {args: args{r: '&'}, want: false},
		"'":     {args: args{r: '\''}, want: false},
		"(":     {args: args{r: '('}, want: false},
		")":     {args: args{r: ')'}, want: false},
		"+":     {args: args{r: '+'}, want: false},
		",":     {args: args{r: ','}, want: false},
		"-":     {args: args{r: '-'}, want: false},
		".":     {args: args{r: '.'}, want: false},
		"0":     {args: args{r: '0'}, want: false},
		"1":     {args: args{r: '1'}, want: false},
		"2":     {args: args{r: '2'}, want: false},
		"3":     {args: args{r: '3'}, want: false},
		"4":     {args: args{r: '4'}, want: false},
		"5":     {args: args{r: '5'}, want: false},
		"6":     {args: args{r: '6'}, want: false},
		"7":     {args: args{r: '7'}, want: false},
		"8":     {args: args{r: '8'}, want: false},
		"9":     {args: args{r: '9'}, want: false},
		";":     {args: args{r: ';'}, want: false},
		"A":     {args: args{r: 'A'}, want: false},
		"B":     {args: args{r: 'B'}, want: false},
		"C":     {args: args{r: 'C'}, want: false},
		"D":     {args: args{r: 'D'}, want: false},
		"E":     {args: args{r: 'E'}, want: false},
		"F":     {args: args{r: 'F'}, want: false},
		"G":     {args: args{r: 'G'}, want: false},
		"H":     {args: args{r: 'H'}, want: false},
		"I":     {args: args{r: 'I'}, want: false},
		"J":     {args: args{r: 'J'}, want: false},
		"K":     {args: args{r: 'K'}, want: false},
		"L":     {args: args{r: 'L'}, want: false},
		"M":     {args: args{r: 'M'}, want: false},
		"N":     {args: args{r: 'N'}, want: false},
		"O":     {args: args{r: 'O'}, want: false},
		"P":     {args: args{r: 'P'}, want: false},
		"Q":     {args: args{r: 'Q'}, want: false},
		"R":     {args: args{r: 'R'}, want: false},
		"S":     {args: args{r: 'S'}, want: false},
		"T":     {args: args{r: 'T'}, want: false},
		"U":     {args: args{r: 'U'}, want: false},
		"V":     {args: args{r: 'V'}, want: false},
		"W":     {args: args{r: 'W'}, want: false},
		"X":     {args: args{r: 'X'}, want: false},
		"Y":     {args: args{r: 'Y'}, want: false},
		"Z":     {args: args{r: 'Z'}, want: false},
		"[":     {args: args{r: '['}, want: false},
		"]":     {args: args{r: ']'}, want: false},
		"_":     {args: args{r: '_'}, want: false},
		"a":     {args: args{r: 'a'}, want: false},
		"b":     {args: args{r: 'b'}, want: false},
		"c":     {args: args{r: 'c'}, want: false},
		"d":     {args: args{r: 'd'}, want: false},
		"e":     {args: args{r: 'e'}, want: false},
		"f":     {args: args{r: 'f'}, want: false},
		"g":     {args: args{r: 'g'}, want: false},
		"h":     {args: args{r: 'h'}, want: false},
		"i":     {args: args{r: 'i'}, want: false},
		"j":     {args: args{r: 'j'}, want: false},
		"k":     {args: args{r: 'k'}, want: false},
		"l":     {args: args{r: 'l'}, want: false},
		"m":     {args: args{r: 'm'}, want: false},
		"n":     {args: args{r: 'n'}, want: false},
		"o":     {args: args{r: 'o'}, want: false},
		"p":     {args: args{r: 'p'}, want: false},
		"q":     {args: args{r: 'q'}, want: false},
		"r":     {args: args{r: 'r'}, want: false},
		"s":     {args: args{r: 's'}, want: false},
		"space": {args: args{r: ' '}, want: false},
		"t":     {args: args{r: 't'}, want: false},
		"u":     {args: args{r: 'u'}, want: false},
		"v":     {args: args{r: 'v'}, want: false},
		"w":     {args: args{r: 'w'}, want: false},
		"x":     {args: args{r: 'x'}, want: false},
		"y":     {args: args{r: 'y'}, want: false},
		"z":     {args: args{r: 'z'}, want: false},
		"Á":     {args: args{r: 'Á'}, want: false},
		"È":     {args: args{r: 'È'}, want: false},
		"É":     {args: args{r: 'É'}, want: false},
		"Ô":     {args: args{r: 'Ô'}, want: false},
		"à":     {args: args{r: 'à'}, want: false},
		"á":     {args: args{r: 'á'}, want: false},
		"ã":     {args: args{r: 'ã'}, want: false},
		"ä":     {args: args{r: 'ä'}, want: false},
		"å":     {args: args{r: 'å'}, want: false},
		"ç":     {args: args{r: 'ç'}, want: false},
		"è":     {args: args{r: 'è'}, want: false},
		"é":     {args: args{r: 'é'}, want: false},
		"ê":     {args: args{r: 'ê'}, want: false},
		"ë":     {args: args{r: 'ë'}, want: false},
		"í":     {args: args{r: 'í'}, want: false},
		"î":     {args: args{r: 'î'}, want: false},
		"ï":     {args: args{r: 'ï'}, want: false},
		"ñ":     {args: args{r: 'ñ'}, want: false},
		"ò":     {args: args{r: 'ò'}, want: false},
		"ó":     {args: args{r: 'ó'}, want: false},
		"ô":     {args: args{r: 'ô'}, want: false},
		"ö":     {args: args{r: 'ö'}, want: false},
		"ø":     {args: args{r: 'ø'}, want: false},
		"ù":     {args: args{r: 'ù'}, want: false},
		"ú":     {args: args{r: 'ú'}, want: false},
		"ü":     {args: args{r: 'ü'}, want: false},
		"ř":     {args: args{r: 'ř'}, want: false},
		"…":     {args: args{r: '…'}, want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := files.IsIllegalRuneForFileNames(tt.args.r); got != tt.want {
				t.Errorf("%s = %v, want %v", fnName, got, tt.want)
			}
		})
	}
}
