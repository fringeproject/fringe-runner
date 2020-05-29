package commands

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fringeproject/fringe-runner/common"
	"github.com/fringeproject/fringe-runner/modules"
	"github.com/fringeproject/fringe-runner/network"
	"github.com/fringeproject/fringe-runner/session"
)

type ServiceCommand struct {
	context  *cli.Context
	session  *session.Session
	config   *common.FringeConfig
	client   common.RunnerClient
	jobsChan chan *common.Job
}

func (s *ServiceCommand) createWorker(id int) {
	logrus.Infof("Create a new worker with id %d", id)

	// The worker wait for an other job in this loop
	for job := range s.jobsChan {
		logrus.Debugf("[Worker %d] Processing job {%s}", id, job.ID)

		// Get the module from the slug
		module, err := s.session.Module(job.Module)
		if err != nil {
			logrus.Warnf("Cannot find module with slug \"%s\"", job.Module)
			logrus.Debug(err)
			continue
		}

		// Create a module context for the execution
		ctx, err := common.NewModuleContext(job.Asset, s.config)
		if err != nil {
			logrus.Warn("Cannot crate module context.")
			logrus.Debug(err)
			continue
		}

		// Run the module
		err = module.Run(ctx)
		if err != nil {
			logrus.Warn("Module execution return an error.")
			logrus.Debug(err)
			continue
		}

		// Get our new shinning assets
		newAssets := ctx.NewAssets

		// Add a retry counter if we've got some problems to upload the results
		maxRetryCount := 3
		curRetryCount := 1
		for curRetryCount < maxRetryCount {
			err = s.client.UpdateJob(job, newAssets)
			if err == nil {
				break
			}

			curRetryCount++
			logrus.Infof("[Worker %d] Retry (%d) to send job update {%s}", id, curRetryCount, job.ID)
			time.Sleep(time.Duration(curRetryCount*10) * time.Second)
		}
	}
}

// Get the next job from the coordinator and add it to the queue
func (s *ServiceCommand) fetchNextJob() error {
	job, err := s.client.RequestJob()

	if err != nil {
		return err
	}

	// There is nothing to do, then pause the runner few seconds
	if job.Module == "" {
		time.Sleep(10 * time.Second)
		return nil
	}

	s.jobsChan <- job

	return nil
}

func (s *ServiceCommand) getRunnerClient() (common.RunnerClient, error) {
	opt := &common.HTTPOptions{
		Headers:        []common.HTTPHeader{},
		Timeout:        time.Second * 20,
		FollowRedirect: true,
		Proxy:          s.config.Proxy,
		VerifyCert:     s.config.VerifyCert,
	}

	coordinator := s.config.FringeCoordinator
	perimeter := s.config.FringePerimeter
	runnerID := s.config.FringeRunnerId
	ruunerToken := s.config.FringeRunnerToken

	client, err := network.NewFringeClient(coordinator, runnerID, ruunerToken, perimeter, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Start an infinite loop to request jobs to the coordinator and execute them
func (s *ServiceCommand) Execute(c *cli.Context, config *common.FringeConfig) error {
	// Create a new session that hold the modules
	sess, err := session.NewSession()
	if err != nil {
		logrus.Warn(err)
		os.Exit(1)
	}
	defer sess.Close()

	// Add the context, session and config to the current command for re-use
	s.context = c
	s.session = sess
	s.config = config

	// Load Fringe modules in the session
	modules.LoadModules(s.session)

	// Get the client to interact with the API
	client, err := s.getRunnerClient()
	if err != nil {
		logrus.Debug(err)
		return fmt.Errorf("Cannot instanciate the Fringe client.")
	}
	s.client = client

	// Notify the server with the list of module available for this runner
	modules, err := s.session.GetModules()
	if err != nil {
		logrus.Debug(err)
		return fmt.Errorf("Cannot generate the module list.")
	}

	err = s.client.SendModuleList(modules)
	if err != nil {
		logrus.Debug(err)
		return fmt.Errorf("Cannot send the module list to the coordinator.")
	}

	// Create runner workers
	// Use the number of CPU to create the workers
	maxWorkers := runtime.NumCPU()
	maxQueueSize := maxWorkers + 1
	s.jobsChan = make(chan *common.Job, maxQueueSize)

	// Create as many worker as we need and set them an ID
	for i := 0; i < maxWorkers; i++ {
		go s.createWorker(i + 1)
	}

	// Infinite loop the process new jobs
	for {
		err := s.fetchNextJob()

		// Something wrong happend, then wait few seconds before continuing
		if err != nil {
			logrus.Error("Something wrong happend while fetching the next job.")
			logrus.Debug(err)
			time.Sleep(10 * time.Second)
		}
	}
}

func newServiceCommand() *ServiceCommand {
	return &ServiceCommand{}
}

func init() {
	common.RegisterCommand("run", "Starts the runner and interacts with the FringeProject API", newServiceCommand(), []cli.Flag{})
}
