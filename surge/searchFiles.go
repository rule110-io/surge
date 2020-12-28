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
	Result []LocalFileListing
	Count  int
}

//SearchRemoteFile runs a paged query
func SearchRemoteFile(Query string, OrderBy string, IsDesc bool, Skip int, Take int) SearchQueryResult {

	var results []FileListing

	ListedFilesLock.Lock()
	for _, file := range ListedFiles {
		if strings.Contains(strings.ToLower(file.FileName), strings.ToLower(Query)) || strings.Contains(strings.ToLower(file.FileHash), strings.ToLower(Query)) {

			result := FileListing{
				FileName:    file.FileName,
				FileHash:    file.FileHash,
				FileSize:    file.FileSize,
				Seeders:     file.seeders,
				NumChunks:   file.NumChunks,
				SeederCount: len(file.seeders),
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
	resultListings := []LocalFileListing{}

	ListedFilesLock.Lock()
	for i := 0; i < len(resultFiles); i++ {
		listing := LocalFileListing{
			ChunkMap:      resultFiles[i].ChunkMap,
			ChunksShared:  resultFiles[i].ChunksShared,
			FileHash:      resultFiles[i].FileHash,
			FileName:      resultFiles[i].FileName,
			FileSize:      resultFiles[i].FileSize,
			IsDownloading: resultFiles[i].IsDownloading,
			IsHashing:     resultFiles[i].IsHashing,
			IsMissing:     resultFiles[i].IsMissing,
			IsPaused:      resultFiles[i].IsPaused,
			IsUploading:   resultFiles[i].IsUploading,
			NumChunks:     resultFiles[i].NumChunks,
			Path:          resultFiles[i].Path,
		}

		if listing.IsUploading {
			listing.Seeders = []string{GetMyAddress()}
		} else {
			listing.Seeders = []string{}
		}

		for _, file := range ListedFiles {
			if file.FileHash == listing.FileHash {
				listing.Seeders = append(listing.Seeders, file.seeders...)
				break
			}
		}
		listing.SeederCount = len(listing.Seeders)
		resultListings = append(resultListings, listing)

		//If file is downloading set progress
		if listing.IsDownloading {
			numChunksLocal := chunksDownloaded(listing.ChunkMap, listing.NumChunks)
			listing.Progress = float32(float64(numChunksLocal) / float64(listing.NumChunks))
		}
	}
	ListedFilesLock.Unlock()

	return LocalFilePageResult{
		Result: resultListings,
		Count:  totalNum,
	}
}
