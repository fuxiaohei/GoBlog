
package utils

import "sort"

type intSorter [][2]int

func (is intSorter) Len() int {
	return len(is)
}

func (is intSorter) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func (is intSorter) Less(i,j int)bool{
	return is[i][1] > is[j][1]
}

// Sort [][2]int by the second element of [2].
func SortInt(s [][2]int){
	sort.Sort(intSorter(s))
}
