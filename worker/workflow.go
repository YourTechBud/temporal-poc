package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"
)

func ExampleWorkflowDefinition(ctx workflow.Context, resource string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 3 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	status := "pending"
	requestID := 0

	workflow.SetQueryHandler(ctx, "status", func() (string, error) {
		return status, nil
	})

	workflow.SetQueryHandler(ctx, "request-id", func() (int, error) {
		return requestID, nil
	})

	makerSignalChan := workflow.GetSignalChannel(ctx, "maker-signal")
	checkerSignalChan := workflow.GetSignalChannel(ctx, "checker-signal")

	for {
		// Create a new draft
		status = "draft"

		fmt.Println("##################### Storing stuff")
		if err := workflow.ExecuteActivity(ctx, TemporaryPushResource, resource).Get(ctx, nil); err != nil {
			return err
		}
		if err := workflow.ExecuteActivity(ctx, StoreResource, resource).Get(ctx, nil); err != nil {
			return err
		}

		// Kick start the ci Pipeline
		fmt.Println("##################### Starting CI Pipeline")

		workflow.ExecuteActivity(ctx, StartResourceCIPipeline, resource)

		for {
			// Get input from developer (send-for-approval | discard | new-push)
			var makerSignal string

			fmt.Println("##################### Waiting for maker input")
			for {
				makerSignalChan.Receive(ctx, &makerSignal)
				arr := strings.Split(makerSignal, ",")
				if arr[0] == strconv.Itoa(requestID) {
					requestID++
					makerSignal = arr[1]
					break
				}
			}

			if makerSignal == "new-push" {
				fmt.Println("##################### New Push")

				// Update the resource object and start process again
				break // break will take us to the above loop
			}

			if makerSignal == "discard" {
				fmt.Println("##################### New Discard")

				// Perform clean up. Maybe cancel ci pipeline, delete temp branch
				status = "discarded"
				return nil
			}

			fmt.Println("##################### New send for approval")

			status = "in-approval"

			// Send for approval just happened. Lets continue
			fmt.Println("##################### Waiting for checker input")
			var checkerSignal string
			for {
				checkerSignalChan.Receive(ctx, &checkerSignal)
				arr := strings.Split(checkerSignal, ",")
				if arr[0] == strconv.Itoa(requestID) {
					requestID++
					checkerSignal = arr[1]
					break
				}
			}

			if checkerSignal == "send-for-review" {
				fmt.Println("##################### New send for review")
				status = "draft"
				continue
			}

			if checkerSignal == "reject" {
				fmt.Println("##################### New reject")

				// Perform clean up. Maybe cancel ci pipeline, delete temp branch
				status = "rejected"
				return nil
			}
			fmt.Println("##################### New approved")

			status = "approved"
			workflow.ExecuteActivity(ctx, PushResource, resource).Get(ctx, nil)

			// Its approved. Run CD Pipeline
			// Maybe we can make cd completed / failed event as a signal
			if err := workflow.ExecuteActivity(ctx, StartResourceCDPipeline, resource).Get(ctx, nil); err != nil {
				// Mark resource as a failed deployment
				return nil
			}

			// Mark resource as completed
			return nil
		}
	}
}
