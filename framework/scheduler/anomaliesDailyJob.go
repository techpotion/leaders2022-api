package scheduler

import (
	"context"
	"time"

	"gitlab.com/techpotion/leadershack2022/api/usecase"
	"gitlab.com/techpotion/leadershack2022/api/utils"
	"go.uber.org/zap"
)

type anomaliesDailyJob struct {
	jobName                            string
	timeSpec                           string
	enrichRequestsWithAnomaliesUsecase *usecase.EnrichRequestsWithAnomaliesUsecase
}

func NewAnomaliesDailyJob(
	timeSpec string,
	enrichRequestsWithAnomaliesUsecase *usecase.EnrichRequestsWithAnomaliesUsecase,
) Job {
	return &anomaliesDailyJob{
		jobName:                            "AnomaliesDailyJob",
		timeSpec:                           timeSpec,
		enrichRequestsWithAnomaliesUsecase: enrichRequestsWithAnomaliesUsecase,
	}
}

func (job *anomaliesDailyJob) GetName() string {
	return job.jobName
}

func (job *anomaliesDailyJob) GetTimeSpec() string {
	return job.timeSpec
}

const maxRetries = 5

func (job *anomaliesDailyJob) Run() {
	z := zap.S().With("context", job.GetName())

	z.Info("Running a scheduled job")

	z.Info("Running merchant_assignments sync")

	for retryCounter := 0; retryCounter <= maxRetries; retryCounter++ {
		if retryCounter == maxRetries {
			z.Errorw("Failed to start merchant assignments sync, max retry amount reached", "retry_counter", maxRetries)
			break
		}

		if jobIsActive := job.enrichRequestsWithAnomaliesUsecase.Execute(context.Background()); jobIsActive {
			z.Infow("Job is already active, retrying...", "table_name", "merchant_assignments", "retry_counter", retryCounter)
			time.Sleep(time.Minute * time.Duration(
				utils.FibonacciRecursion(retryCounter+1),
			))

			continue
		}

		break
	}

	z.Info("Finished merchant_assignments sync")
}
