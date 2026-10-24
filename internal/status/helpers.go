package status

import "errors"

func StorePausedWorker(store *Store, id, name string) error {
	if store == nil {
		return errors.New("cannot store paused worker status: store is nil")
	}

	workerStatus, err := NewWorkerStatus(
		id,
		name,
		Paused,
		[]MonitoredService{},
		"Worker is disabled by configuration",
	)
	if err != nil {
		return err
	}

	store.Set(workerStatus)
	return nil
}
