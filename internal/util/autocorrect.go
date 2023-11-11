package util

import "github.com/lithammer/fuzzysearch/fuzzy"

func LevensteinDistance(s1, s2 string) int {
	// degenerate cases
	if s1 == s2 {
		return 0
	}

	l1 := len(s1)
	l2 := len(s2)
	// convert string to rune array
	r1 := []rune(s1)
	r2 := []rune(s2)

	// create two-dimensional array
	matrix := make([][]int, l1+1)
	for i := 0; i < l1+1; i++ {
		matrix[i] = make([]int, l2+1)
	}

	// initialize matrix
	for i := 0; i < l1+1; i++ {
		matrix[i][0] = i
	}

	for j := 0; j < l2+1; j++ {
		matrix[0][j] = j
	}

	// calculate matrix
	for i := 1; i < l1+1; i++ {
		for j := 1; j < l2+1; j++ {
			if r1[i-1] == r2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = matrixMin(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+1)
			}
		}
	}

	return matrix[l1][l2]
}

func matrixMin(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func FuzzyFind(arg string, choices []string) string {
	lowest := 1<<31 - 1
	best := ""
	f := fuzzy.RankFindFold(arg, choices)
	if len(f) == 0 {
		for i := range choices {
			lv := LevensteinDistance(arg, choices[i])
			if lv < lowest {
				lowest = lv
				best = choices[i]
			}
		}
		return best
	}

	for i := range f {
		if f[i].Distance < lowest {
			lowest = f[i].Distance
			best = f[i].Target
		}
	}
	return best
}
