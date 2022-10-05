package main

import (
	"context"
	"fmt"
	"time"
)

func TemporaryPushResource(ctx context.Context, resource string) error {
	fmt.Printf("Pushed '%s' to temp branch in git\n", resource)
	return nil
}

func PushResource(ctx context.Context, resource string) error {
	fmt.Printf("Pushed '%s' to main branch in git\n", resource)
	return nil
}

func StoreResource(ctx context.Context, resource string) error {
	fmt.Printf("Stored '%s' to postgres\n", resource)
	return nil
}

func StartResourceCIPipeline(ctx context.Context, resource string) error {
	fmt.Printf("Started CI pipeline for '%s'\n", resource)
	time.Sleep(30 * time.Second)
	fmt.Printf("Finished CI pipeline for '%s'\n", resource)
	return nil
}

func StartResourceCDPipeline(ctx context.Context, resource string) error {
	// Make this idempotent. Only start pipeline if it isn't running already
	fmt.Printf("Started CD pipeline for '%s'\n", resource)
	time.Sleep(30 * time.Second)
	fmt.Printf("Finished CD pipeline for '%s'\n", resource)
	return nil
}
