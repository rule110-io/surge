package surge

import (
	"sort"
	"strings"
)

//SearchQueryResult is a paging query result for file searches
type SearchQueryResult struct {
	Result []FileListing
	Count  int
}

//LocalFilePageResult is a paging query result for tracked files
type LocalFilePageResult struct {
	Result []File
	Count  int
}

//SearchRemoteFile runs a paged query
func SearchRemoteFile(Query string, OrderBy string, IsDesc bool, Skip int, Take int) SearchQueryResult {
	defer RecoverAndLog()
	var results []FileListing

	ListedFilesLock.Lock()
	for _, file := range ListedFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) {

			result := FileListing{
				FileName:    file.FileName,
				FileHash:    file.FileHash,
				FileSize:    file.FileSize,
				Seeders:     file.Seeders,
				NumChunks:   file.NumChunks,
				SeederCount: len(file.Seeders),
			}

			tracked, err := dbGetFile(result.FileHash)

			//only add non-local files to the result
			if err != nil && tracked == nil {
				results = append(results, result)
			}

		}
	}
	ListedFilesLock.Unlock()

	switch OrderBy {
	case "FileName":
		if !IsDesc {
			sort.Sort(sortByFileNameAsc(results))
		} else {
			sort.Sort(sortByFileNameDesc(results))
		}
	case "FileSize":
		if !IsDesc {
			sort.Sort(sortByFileSizeAsc(results))
		} else {
			sort.Sort(sortByFileSizeDesc(results))
		}
	default:
		if !IsDesc {
			sort.Sort(sortBySeederCountAsc(results))
		} else {
			sort.Sort(sortBySeederCountDesc(results))
		}
	}

	left := Skip
	right := Skip + Take

	if left > len(results) {
		left = len(results)
	}

	if right > len(results) {
		right = len(results)
	}

	return SearchQueryResult{
		Result: results[left:right],
		Count:  len(results),
	}
}

//SearchLocalFile runs a paged query
func SearchLocalFile(Query string, OrderBy string, IsDesc bool, Skip int, Take int) LocalFilePageResult {
	defer RecoverAndLog()
	var results []FileListing

	resultFiles := []File{}

	allFiles := dbGetAllFiles()
	for _, file := range allFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) {
			resultFiles = append(resultFiles, file)
		}
	}

	totalNum := len(resultFiles)
	for i := 0; i < len(resultFiles); i++ {
		resultFiles[i].ChunkMap = nil
	}

	switch OrderBy {
	case "FileName":
		if !IsDesc {
			sort.Sort(sortByFileNameAsc(results))
		} else {
			sort.Sort(sortByFileNameDesc(results))
		}
	case "FileSize":
		if !IsDesc {
			sort.Sort(sortByFileSizeAsc(results))
		} else {
			sort.Sort(sortByFileSizeDesc(results))
		}
	default:
		if !IsDesc {
			sort.Sort(sortBySeederCountAsc(results))
		} else {
			sort.Sort(sortBySeederCountDesc(results))
		}
	}

	left := Skip
	right := Skip + Take

	if left > len(resultFiles) {
		left = len(resultFiles)
	}

	if right > len(resultFiles) {
		right = len(resultFiles)
	}

	//Subset
	resultFiles = resultFiles[left:right]

	for i := 0; i < len(resultFiles); i++ {
		ListedFilesLock.Lock()

		for _, file := range ListedFiles {
			resultFiles[i].Seeders = []string{GetMyAddress()}
			if file.FileHash == resultFiles[i].FileHash {
				resultFiles[i].Seeders = file.Seeders
				resultFiles[i].Seeders = append(resultFiles[i].Seeders, GetMyAddress())
				resultFiles[i].SeederCount = len(resultFiles[i].Seeders)
				break
			}
		}

		if len(resultFiles[i].Seeders) == 0 && (resultFiles[i].IsUploading || resultFiles[i].IsHashing) {
			resultFiles[i].Seeders = []string{GetMyAddress()}
			resultFiles[i].SeederCount = len(resultFiles[i].Seeders)
		}

		ListedFilesLock.Unlock()
	}

	return LocalFilePageResult{
		Result: resultFiles,
		Count:  totalNum,
	}
}
