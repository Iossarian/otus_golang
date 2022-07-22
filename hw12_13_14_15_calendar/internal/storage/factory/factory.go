package factory

import (
	"errors"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/config"
	memoryStorage "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlStorage "github.com/Iossarian/otus_golang/hw12_13_14_15_calendar/internal/storage/sql"
)

const memoryStorageName = "memory"
const sqlStorageName = "sql"

var ErrUnknownStorageType = errors.New("can not get storage source")

func GetStorage(c config.Config) (app.Storage, error) {
	switch c.StorageSource {
	case memoryStorageName:
		return memoryStorage.New(), nil
	case sqlStorageName:
		return sqlStorage.New(c), nil
	default:
		return nil, ErrUnknownStorageType
	}
}
