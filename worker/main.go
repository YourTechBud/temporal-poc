package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	client, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to intialise temporal client", err)
	}
	defer client.Close()

	w := worker.New(client, "tq", worker.Options{})
	w.RegisterWorkflow(ExampleWorkflowDefinition)

	w.RegisterActivity(TemporaryPushResource)
	w.RegisterActivity(PushResource)
	w.RegisterActivity(StoreResource)
	w.RegisterActivity(StartResourceCIPipeline)
	w.RegisterActivity(StartResourceCDPipeline)

	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
