package controller

import (
	coreV1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubeclientset "k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/jiaxuanzhou/jx-controller/pkg/apis/jx/v1alpha1"
	clientset "github.com/jiaxuanzhou/jx-controller/pkg/client/clientset/versioned"
	"github.com/jiaxuanzhou/jx-controller/pkg/client/clientset/versioned/scheme"
	jxInformersV1alpha1 "github.com/jiaxuanzhou/jx-controller/pkg/client/informers/externalversions/jx/v1alpha1"
	listers "github.com/jiaxuanzhou/jx-controller/pkg/client/listers/jx/v1alpha1"
	"github.com/prometheus/prometheus/common/log"

	"fmt"
	"time"
)

const controllerName = "jx-controller"

// Controller is the controller implementation for Foo resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// jxTaskclientset is a clientset for our own API group
	jxTaskclientset clientset.Interface

	// podLister can list/get pods from the shared informer's store.
	podLister corelisters.PodLister
	// podInformerSynced returns true if the pod store has been synced at least once.
	podInformerSynced cache.InformerSynced

	jxTaskLister  listers.JxTaskLister
	jxTasksSynced cache.InformerSynced

	// jxTaskInformer is a temporary field for unstructured informer support.
	jxTaskInformer cache.SharedIndexInformer
	// tfJobInformerSynced returns true if the tfjob store has been synced at least once.
	jxTaskInformerSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewJxController returns a new JxTask controller.
func NewJxController(
	kubeClientSet kubeclientset.Interface,
	jxTaskclientset clientset.Interface,
	jxInformers jxInformersV1alpha1.JxTaskInformer,
	kubeInformerFactory kubeinformers.SharedInformerFactory) *Controller {

	scheme.AddToScheme(scheme.Scheme)

	log.Debug("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(log.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, coreV1.EventSource{Component: controllerName})

	// Create new TFJobController.
	jc := &Controller{
		kubeclientset:   kubeClientSet,
		jxTaskclientset: jxTaskclientset,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), v1alpha1.CRDKindPlural),
		recorder:        recorder,
	}

	// Set up an event handler for when tfjob resources change.
	jxInformers.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    jc.addJxTask,
		UpdateFunc: jc.updateJxTask,
		// This will enter the sync loop and no-op,
		// because the tfjob has been deleted from the store.
		DeleteFunc: jc.enqueueJxTask,
	})

	jc.jxTaskInformer = jxInformers.Informer()
	jc.jxTaskLister = jxInformers.Lister()
	jc.jxTaskInformerSynced = jxInformers.Informer().HasSynced

	// Create pod informer.
	podInformer := kubeInformerFactory.Core().V1().Pods()

	// Set up an event handler for when pod resources change
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    jc.addPod,
		UpdateFunc: jc.updatePod,
		DeleteFunc: jc.deletePod,
	})

	jc.podLister = podInformer.Lister()
	jc.podInformerSynced = podInformer.Informer().HasSynced

	return jc
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (jc *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer jc.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches.
	log.Info("Starting TFJob controller")

	// Wait for the caches to be synced before starting workers.
	log.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, jc.jxTaskInformerSynced); !ok {
		return fmt.Errorf("failed to wait for jxtask caches to sync")
	}

	if ok := cache.WaitForCacheSync(stopCh, jc.podInformerSynced); !ok {
		return fmt.Errorf("failed to wait for pod caches to sync")
	}

	log.Infof("Starting %v workers", threadiness)
	// Launch workers to process TFJob resources.
	for i := 0; i < threadiness; i++ {
		go wait.Until(jc.runWorker, time.Second, stopCh)
	}

	log.Info("Started workers")
	<-stopCh
	log.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (jc *Controller) runWorker() {
	for jc.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (jc *Controller) processNextWorkItem() bool {
	return true
}
