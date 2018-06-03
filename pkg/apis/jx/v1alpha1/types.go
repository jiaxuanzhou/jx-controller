package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CRDKind       = "jxtask"
	CRDKindPlural = "jxtasks"
	CRDGroup      = "jiaxuan.org"
	CRDVersion    = "v1alpha1"
	// Value of the APP label that gets applied to a lot of entities.
	AppLabel = "by-jiaxuan"
	// Defaults for the Spec
	JXPort   = 9999
	Replicas = 1
	EnvJxNamespace = "JX_NAMESPACE"
)

const (
	DefaultJXContainer string = "busybox"
	DefaultImage       string = "busybox"
)

// TaskType determines how a set of tasks are handled. Support three kinds of tasks: online/offline/batch.
type TaskType string

const (
	Batch   TaskType = "BATCH"
	OnLine  TaskType = "ONLINE"
	OffLine TaskType = "OFFLINE"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jxtask

// JxTask describes jxtask info
type JxTask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              JxSpec       `json:"spec"`
	Status            JxTaskStatus `json:"status"`
}

// JxSpec structure for storing the jxTask specifications
type JxSpec struct {
	RuntimeId string

	// Tasks specifies the jx task to run.
	Tasks []*Task `json:"task"`

	// Jx defines default image used for jx tasks
	//JxImage string `json:"jxImage,omitempty"`

	// TerminationPolicy specifies the condition that the jxtask should be considered finished.
	// todo: or can consider finalizer.
	TerminationPolicy *TerminationPolicySpec `json:"terminationPolicy,omitempty"`

	// SchedulerName specifies the name of scheduler which should handle the tasks
	SchedulerName string `json:"schedulerName,omitempty"`
}

// TerminationPolicySpec structure for storing specifications for process termination
type TerminationPolicySpec struct {
	// Chief policy waits for a particular task (which is the chief) to exit.
	Chief *ChiefSpec `json:"chief,omitempty"`
}

// ChiefSpec structure storing the task name and task index
type ChiefSpec struct {
	TaskName  string `json:"taskName"`
	TaskIndex int    `json:"taskIndex"`
}

// Task might be useful if you wanted to have a separate set of workers to do eval.
type Task struct {
	// Name is the task of one specific task
	Name string `json:"name,omitempty"`
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// More info: http://kubernetes.io/docs/user-guide/replication-controller#what-is-a-replication-controller
	// +optional
	Replicas *int32              `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
	Template *v1.PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,3,opt,name=template"`
	// JxPort is the port to use for jx services.
	JxPort   *int32 `json:"jxPort,omitempty" protobuf:"varint,1,opt,name=jxPort"`
	TaskType `json:"jxTaskType"`
}

// JxTaskPhase is a enum to store the phase of jx task
type JxTaskPhase string

const (
	JxTaskPhaseNone     JxTaskPhase = ""
	JxTaskPhaseCreating JxTaskPhase = "Creating"
	JxTaskPhaseRunning  JxTaskPhase = "Running"
	JxTaskPhaseCleanUp  JxTaskPhase = "CleanUp"
	JxTaskPhaseFailed   JxTaskPhase = "Failed"
	JxTaskPhaseDone     JxTaskPhase = "Done"
)

// State is a enum to store the state of  jx task
type State string

const (
	StateUnknown   State = "Unknown"
	StateRunning   State = "Running"
	StateSucceeded State = "Succeeded"
	StateFailed    State = "Failed"
)

// JxTaskStatus is a structure for storing the status of jx tasks
type JxTaskStatus struct {
	// Phase is the Jx task running phase
	Phase  JxTaskPhase `json:"phase"`
	Reason string      `json:"reason"`

	// State indicates the state of the job.
	State State `json:"state"`

	// TaskStatuses specifies the status of each JX task.
	TaskStatuses []*TaskStatus `json:"taskStatuses"`
}

// TaskState is a enum to store the status of task
type TaskState string

const (
	TaskStateUnknown   TaskState = "Unknown"
	TaskStateRunning   TaskState = "Running"
	TaskStateFailed    TaskState = "Failed"
	TaskStateSucceeded TaskState = "Succeeded"
)

// TaskStatus  is a structure for storing the status of task
type TaskStatus struct {
	TaskType `json:"jx_task_type"`

	// State is the overall state of the task
	State TaskState `json:"state"`

	// TasksStates provides the number of task in each status.
	TasksStates map[TaskState]int
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jxtasks

// JxTaskList is a list of JxTask clusters.
type JxTaskList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items is a list of JxTasks
	Items []JxTask `json:"items"`
}

// ControllerConfig is a structure for storing the controller configuration
type ControllerConfig struct {
	// Accelerators is a map from the name of the accelerator to the config for that accelerator.
	// This should match the value specified as a container limit.
	// e.g. alpha.kubernetes.io/nvidia-gpu
	Accelerators map[string]AcceleratorConfig

	// Path to the file containing the grpc server source
	GrpcServerFilePath string
}

// AcceleratorConfig represents accelerator volumes to be mounted into container along with environment variables.
type AcceleratorConfig struct {
	Volumes []AcceleratorVolume
	EnvVars []EnvironmentVariableConfig
}

// AcceleratorVolume represents a host path that must be mounted into
// each container that needs to use GPUs.
type AcceleratorVolume struct {
	Name      string
	HostPath  string
	MountPath string
}

// EnvironmentVariableConfig represents the environment variables and their values.
type EnvironmentVariableConfig struct {
	Name  string
	Value string
}
