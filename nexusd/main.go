package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Supernomad/nexus/nexusd/common"
	"github.com/Supernomad/nexus/nexusd/filter"
	"github.com/Supernomad/nexus/nexusd/iface"
	"github.com/Supernomad/nexus/nexusd/worker"
)

func main() {
	log := common.NewLogger()
	cfg, err := common.NewConfig(log)
	if err != nil {
		log.Error.Fatalln("[MAIN]", "Configuration issue encountered:", err)
	}

	backend := iface.New(iface.UdpSocket, log, cfg)
	err = backend.Open()
	if err != nil {
		log.Error.Fatalln("[MAIN]", "Error bringing up the backend networking interface:", err)
	}

	filters := make([]filter.Filter, 0)
	worker := worker.New(log, cfg, backend, filters)
	for i := 0; i < cfg.NumWorkers; i++ {
		worker.Start(i)
	}

	signals := make(chan os.Signal, 1)
	done := make(chan struct{}, 1)

	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-signals
		switch {
		case sig == syscall.SIGHUP:
			log.Info.Println("[MAIN]", "Recieved reload signal from user. Reloading process.")

			backendFDs := backend.GetFDs()

			files := make([]uintptr, 3+cfg.NumWorkers*2)
			files[0] = os.Stdin.Fd()
			files[1] = os.Stdout.Fd()
			files[2] = os.Stderr.Fd()

			for i := 0; i < cfg.NumWorkers; i++ {
				files[3+i] = uintptr(backendFDs[i])
			}

			os.Setenv(common.RollingRestart, "restart triggered")
			env := os.Environ()
			attr := &syscall.ProcAttr{
				Env:   env,
				Files: files,
			}

			pid, err := syscall.ForkExec(os.Args[0], os.Args, attr)
			if err != nil {
				log.Error.Fatalln("[MAIN]", "Rolling restart error:", err)
			}

			ioutil.WriteFile(cfg.PidFile, []byte(strconv.Itoa(pid)), os.ModePerm)
			done <- struct{}{}
		case sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGKILL:
			log.Info.Println("[MAIN]", "Recieved termination signal from user. Terminating process.")
			done <- struct{}{}
		}
	}()

	<-done
	worker.Stop()
	backend.Close()
}
