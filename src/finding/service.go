package main

const (
	defaultSortDirection string = "asc"
	defaultLimit         int32  = 200
)

type findingService struct {
	repository findingRepository
}

// TODO remove after fix response interface type change to int64 from uint32
func convertToUint32(v int64) uint32 {
	return uint32(v)
}
