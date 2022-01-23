package surge

import (
	"strings"

	"github.com/rule110-io/surge/backend/models"
)

type sortBySeederCountAsc []models.File

func (a sortBySeederCountAsc) Len() int { return len(a) }
func (a sortBySeederCountAsc) Less(i, j int) bool {
	return len(GetSeeders(a[i].FileHash)) < len(GetSeeders(a[j].FileHash))
}
func (a sortBySeederCountAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortBySeederCountDesc []models.File

func (a sortBySeederCountDesc) Len() int { return len(a) }
func (a sortBySeederCountDesc) Less(i, j int) bool {
	return len(GetSeeders(a[i].FileHash)) > len(GetSeeders(a[j].FileHash))
}
func (a sortBySeederCountDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileNameAsc []models.File

func (a sortByFileNameAsc) Len() int { return len(a) }
func (a sortByFileNameAsc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) < strings.ToLower(a[j].FileName)
}
func (a sortByFileNameAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileNameDesc []models.File

func (a sortByFileNameDesc) Len() int { return len(a) }
func (a sortByFileNameDesc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) > strings.ToLower(a[j].FileName)
}
func (a sortByFileNameDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByFileSizeAsc []models.File

func (a sortByFileSizeAsc) Len() int           { return len(a) }
func (a sortByFileSizeAsc) Less(i, j int) bool { return a[i].FileSize < a[j].FileSize }
func (a sortByFileSizeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortByFileSizeDesc []models.File

func (a sortByFileSizeDesc) Len() int           { return len(a) }
func (a sortByFileSizeDesc) Less(i, j int) bool { return a[i].FileSize > a[j].FileSize }
func (a sortByFileSizeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

/////FileListing

type sortByListingSeederCountAsc []models.FileListing

func (a sortByListingSeederCountAsc) Len() int { return len(a) }
func (a sortByListingSeederCountAsc) Less(i, j int) bool {
	return len(GetSeeders(a[i].FileHash)) < len(GetSeeders(a[j].FileHash))
}
func (a sortByListingSeederCountAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByListingSeederCountDesc []models.FileListing

func (a sortByListingSeederCountDesc) Len() int { return len(a) }
func (a sortByListingSeederCountDesc) Less(i, j int) bool {
	return len(GetSeeders(a[i].FileHash)) > len(GetSeeders(a[j].FileHash))
}
func (a sortByListingSeederCountDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByListingFileNameAsc []models.FileListing

func (a sortByListingFileNameAsc) Len() int { return len(a) }
func (a sortByListingFileNameAsc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) < strings.ToLower(a[j].FileName)
}
func (a sortByListingFileNameAsc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByListingFileNameDesc []models.FileListing

func (a sortByListingFileNameDesc) Len() int { return len(a) }
func (a sortByListingFileNameDesc) Less(i, j int) bool {
	return strings.ToLower(a[i].FileName) > strings.ToLower(a[j].FileName)
}
func (a sortByListingFileNameDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type sortByListingFileSizeAsc []models.FileListing

func (a sortByListingFileSizeAsc) Len() int           { return len(a) }
func (a sortByListingFileSizeAsc) Less(i, j int) bool { return a[i].FileSize < a[j].FileSize }
func (a sortByListingFileSizeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type sortByListingFileSizeDesc []models.FileListing

func (a sortByListingFileSizeDesc) Len() int           { return len(a) }
func (a sortByListingFileSizeDesc) Less(i, j int) bool { return a[i].FileSize > a[j].FileSize }
func (a sortByListingFileSizeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
