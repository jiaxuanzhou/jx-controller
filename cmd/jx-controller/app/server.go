package app

import (
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/jiaxuanzhou/jx-controller/cmd/jx-controller/app/options"
	"github.com/jiaxuanzhou/jx-controller/pkg/apis/jx/v1alpha1"
	jxtaskclientset "github.com/jiaxuanzhou/jx-controller/pkg/client/clientset/versioned"
	"github.com/jiaxuanzhou/jx-controller/pkg/client/clientset/versioned/scheme"
	"github.com/jiaxuanzhou/jx-controller/pkg/controller"
	"github.com/jiaxuanzhou/jx-controller/pkg/util/signals"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeinformers "k8s.io/client-go/informers"
	kubeclientset "k8s.io/client-go/kubernetes"
	restclientset "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	election "k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

const (
	apiVersion = "v1alpha1"
)

var (
	// leader election config
	leaseDuration = 15 * time.Second
	renewDuration = 5 * time.Second
	retryPeriod   = 3 * time.Second
	resyncPeriod  = 30 * time.Second
)

// KubeConfigPathEnv will be import from configMap or just from env of the deployment
const KubeConfigPathEnv = "KUBECONFIG"

func Run(opt *options.ServerOption) error {

	namespace := os.Getenv(v1alpha1.EnvJxNamespace)
	if len(namespace) == 0 {
		glog.Infof("EnvJxNamespace not set, use default namespace")
		namespace = metav1.NamespaceDefault
	}

	// Set up signals so we handle the first shutdown signal gracefully.
	stopCh := signals.SetupSignalHandler()

	// Note: ENV KUBECONFIG will overwrite user defined Kubeconfig option.
	if len(os.Getenv(KubeConfigPathEnv)) > 0 {
		// use the current context in kubeconfig
		// This is very useful for running locally.
		opt.Kubeconfig = os.Getenv(KubeConfigPathEnv)
	}

	// Get kubernetes config.
	kcfg, err := clientcmd.BuildConfigFromFlags(opt.Master, opt.Kubeconfig)
	if err != nil {
		glog.Fatalf("Error building kubeconfig: %s", err.Error())
	}

	// Create kube clients.
	kubeClientSet, leaderElectionClientSet, jxTaskClientSet, err := createClientSets(kcfg)
	if err != nil {
		return err
	}

	// 1, Create informer factory.
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClientSet, resyncPeriod)

	jxInformer := controller.NewJxTaskInformer(kcfg)
	// 2, Create jx controller.
	jc := controller.NewJxController(kubeClientSet, jxTaskClientSet, jxInformer, kubeInformerFactory)

	// 3, Start informer goroutines.
	go kubeInformerFactory.Start(stopCh)

	go jxInformer.Informer().Run(stopCh)
	// Set leader election start function.
	run := func(<-chan struct{}) {
		if err := jc.Run(1, stopCh); err != nil {
			glog.Errorf("Failed to run the controller: %v", err)
		}
	}

	id, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("Failed to get hostname: %v", err)
	}

	// Prepare event clients.
	eventBroadcaster := record.NewBroadcaster()
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: v1alpha1.CRDKindPlural})

	rl := &resourcelock.EndpointsLock{
		EndpointsMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      v1alpha1.CRDKindPlural,
		},
		Client: leaderElectionClientSet.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: recorder,
		},
	}

	// Start leader election.
	election.RunOrDie(election.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: leaseDuration,
		RenewDeadline: renewDuration,
		RetryPeriod:   retryPeriod,
		Callbacks: election.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				glog.Fatalf("leader election lost")
			},
		},
	})

	return nil
}

// fill the funcs when client code generated by v1alpha1
func createClientSets(config *restclientset.Config) (kubeclientset.Interface, kubeclientset.Interface, jxtaskclientset.Interface, error) {
	kubeClientSet, err := kubeclientset.NewForConfig(restclientset.AddUserAgent(config, "jxtask"))
	if err != nil {
		return nil, nil, nil, err
	}

	leaderElectionClientSet, err := kubeclientset.NewForConfig(restclientset.AddUserAgent(config, "leader-election"))
	if err != nil {
		return nil, nil, nil, err
	}

	jxTaskClientSet, err := jxtaskclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	return kubeClientSet, leaderElectionClientSet, jxTaskClientSet, nil
}
