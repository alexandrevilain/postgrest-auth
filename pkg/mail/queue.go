package mail

import (
	"github.com/alexandrevilain/postgrest-auth/pkg/config"
	"github.com/labstack/gommon/log"
)

// EmailSendRequest represents a mail sending request
type EmailSendRequest struct {
	To      string
	Title   string
	Content string
}

// Worker is the struct maintaining worker's state
type Worker struct {
	ID          int
	WorkerQueue chan EmailSendRequest
	QuitChan    chan bool
	sender      *sender
	logger      *log.Logger
}

// NewSenderWorker creates, and returns a new Worker object.
func NewSenderWorker(workerQueue chan EmailSendRequest, config *config.Email, logger *log.Logger) *Worker {
	return &Worker{
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		sender:      newSender(config),
		logger:      logger,
	}
}

// Start launches the worker by starting a goroutine
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case work := <-w.WorkerQueue:
				// Receive a work request.
				w.logger.Info("worker: Received email sending request \n")

				err := w.sender.sendEmail(work.To, work.Title, work.Content)
				if err != nil {
					w.logger.Errorf("An error occured while sending email: %v \n", err.Error())
				}

			case <-w.QuitChan:
				// We have been asked to stop.
				w.logger.Info("Stopping worker ...\n")
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
