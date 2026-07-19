package repository

type CacheRepo interface {
	StoreStrings(key string, strings []string) error
	RetrieveStrings(key string, startIndex int, endIndex int) ([]string, error)
	ClearCache(key string) error
}
