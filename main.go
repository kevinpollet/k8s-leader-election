package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var (
	leaseLockName      = flag.String("lease-lock-name", "leader-election", "The lease lock name.")
	leaseLockNamespace = flag.String("lease-lock-namespace", "default", "The lease lock namespace.")
)

func main() {
	flag.Parse()

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	id := uuid.New()
	kubeClient := kubernetes.NewForConfigOrDie(config)

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      *leaseLockName,
			Namespace: *leaseLockNamespace,
		},
		Client: kubeClient.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id.String(),
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,             // Release the lock when the run context is cancelled.
		LeaseDuration:   60 * time.Second, // The duration that non-leader candidates will wait to force acquire leadership.
		RenewDeadline:   15 * time.Second, // The duration that the acting master will retry refreshing leadership before giving up.
		RetryPeriod:     5 * time.Second,  // The duration the LeaderElector clients should wait between tries of actions.
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(_ context.Context) {
				log.Println("Start leading")
			},
			OnStoppedLeading: func() {
				log.Println("Stop leading")
			},
			OnNewLeader: func(identity string) {
				// We are the leader.
				if identity == id.String() {
					return
				}

				log.Printf("New leader elected: %s\n", identity)
			},
		},
	})
}
