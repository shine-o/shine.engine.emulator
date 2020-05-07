package service

type queryMapEvent struct {
	id int
	zm chan <- *zoneMap
	err chan error
}

func (e *queryMapEvent) erroneous() <- chan error {
	return e.err
}