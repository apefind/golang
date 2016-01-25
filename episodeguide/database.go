package episodeguide

import (
	"fmt"
	"time"
)

type SeriesDataBaseReader interface {
	GetSeries(title string) (*Series, error)
	UpdateSeries(series *Series) error
}

func GetSeriesFromDataBase(title string, readers []SeriesDataBaseReader, timeout time.Duration) (*Series, error) {
	var series *Series
	var err error
	done := make(chan error, 1)
	for _, r := range readers {
		go func(r SeriesDataBaseReader) {
			series, err = r.GetSeries(title)
			if err == nil {
				done <- err
			}
		}(r)
	}
	select {
	case <-time.After(timeout):
		close(done)
		return nil, fmt.Errorf("timeout after %s", timeout)
	case err = <-done:
		close(done)
		return series, err
	}
}
