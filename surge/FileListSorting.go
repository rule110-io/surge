package surge

import "strings"

/*type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}*/

type sortBySeederCountAsc []FileListing

func (a sortBySeederCountAsc) Len() int           { return len(a) }
func (a sortBySeederCountAsc) Less(i, j int) bool { return a[i].SeederCount < a[j].SeederCount }
func (a sortBySeederCountAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortBySeederCountDesc []FileListing

func (a sortBySeederCountDesc) Len() int           { return len(a) }
func (a sortBySeederCountDesc) Less(i, j int) bool { return a[i].SeederCount > a[j].SeederCount }
func (a sortBySeederCountDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortByFileNameAsc []FileListing

func (a sortByFileNameAsc) Len() int { return len(a) }
func (a sortByFileNameAsc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) < strings.ToLower(a[j].FileName)
}
func (a sortByFileNameAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileNameDesc []FileListing

func (a sortByFileNameDesc) Len() int { return len(a) }
func (a sortByFileNameDesc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) > strings.ToLower(a[j].FileName)
}
func (a sortByFileNameDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileSizeAsc []FileListing

func (a sortByFileSizeAsc) Len() int           { return len(a) }
func (a sortByFileSizeAsc) Less(i, j int) bool { return a[i].FileSize < a[j].FileSize }
func (a sortByFileSizeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortByFileSizeDesc []FileListing

func (a sortByFileSizeDesc) Len() int           { return len(a) }
func (a sortByFileSizeDesc) Less(i, j int) bool { return a[i].FileSize > a[j].FileSize }
func (a sortByFileSizeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortLocalByFileNameAsc []File

func (a sortLocalByFileNameAsc) Len() int { return len(a) }
func (a sortLocalByFileNameAsc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) < strings.ToLower(a[j].FileName)
}
func (a sortLocalByFileNameAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortLocalByFileNameDesc []File

func (a sortLocalByFileNameDesc) Len() int { return len(a) }
func (a sortLocalByFileNameDesc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) > strings.ToLower(a[j].FileName)
}
func (a sortLocalByFileNameDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
