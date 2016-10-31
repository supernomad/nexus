package worker

import (
	"github.com/Supernomad/nexus/common"
	"github.com/Supernomad/nexus/filter"
	"github.com/Supernomad/nexus/iface"
)

// Worker object to handle recieving, filtering and sending packets
type Worker struct {
	filters []filter.Filter
	backend iface.Iface
	log     *common.Logger
	cfg     *common.Config
	done    bool
}

func (worker *Worker) pipeline(queue int) error {
	packet, err := worker.backend.Read(queue)
	if err != nil {
		return err
	}

	for i := 0; i < len(worker.filters); i++ {
		if drop := worker.filters[i].Drop(packet); drop {
			return nil
		}
	}

	return worker.backend.Write(queue, packet)
}

// Start the worker object
func (worker *Worker) Start(queue int) {
	go func() {
		for !worker.done {
			if err := worker.pipeline(queue); err != nil {
				worker.log.Error.Println("[WORKER]", "Error during work pipeline:", err)
			}
		}
	}()
}

// Stop the worker object gracefully
func (worker *Worker) Stop() {
	worker.done = true
}

// New work object
func New(log *common.Logger, cfg *common.Config, backend iface.Iface, filters []filter.Filter) *Worker {
	return &Worker{
		filters: filters,
		backend: backend,
		log:     log,
		cfg:     cfg,
		done:    false,
	}
}
