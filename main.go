package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

func main() {
	// create a scheduler
	scheduler, err := gocron.NewScheduler(
		gocron.WithLocation(time.Local),
	)
	if err != nil {
		// handle error
		fmt.Printf("定时器创建失败: %v\n", err)
	}

	// add a job to the scheduler
	_, err = scheduler.NewJob(
		gocron.DurationJob(
			time.Second*10, // 每10s执行  Execute every 10 seconds
		),
		gocron.NewTask(
			func(a string, b int) {
				fmt.Printf("task1 %v,%v\n", a, b)
			},
			"hello",
			1,
		),
	)
	if err != nil {
		// handle error
		fmt.Printf("task1 创建失败: %v\n", err)
	}

	// add a job to the scheduler
	_, err = scheduler.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(6, 0, 0)), // 每天6点执行  Execute at 6:00 every day
		),
		gocron.NewTask(
			func() {
				fmt.Printf("task2\n")
			},
		),
		gocron.WithStartAt(gocron.WithStartImmediately()), // 启动时立即执行一次 Execute once immediately at startup
		// gocron.LimitModeWait 单一并发执行,排队等待  Single job, waiting in queue
		// gocron.LimitModeReschedule 单一并发执行,放弃跳过  Single job, give up skipping
		gocron.WithSingletonMode(gocron.LimitModeWait),
	)
	if err != nil {
		// handle error
		fmt.Printf("task2 创建失败: %v\n", err)
	}

	// add a job to the scheduler
	_, err = scheduler.NewJob(
		gocron.DurationJob(
			time.Second*10,
		),
		gocron.NewTask(
			func() error {
				//fmt.Println("task3 执行成功")
				//return nil
				return errors.New("task3 返回错误")
			},
		),
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					// 当作业返回错误时执行操作  do something when the job returns an error
					fmt.Printf("task3 执行失败,%v,%v,%v\n", jobID, jobName, err)
				},
			),
		),
	)
	if err != nil {
		// handle error
		fmt.Printf("task3 创建失败: %v\n", err)
	}

	scheduler.Start()

	select {} // wait forever
}
