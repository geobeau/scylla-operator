package naming

// These labels are only used on the ClusterIP services
// acting as each member's identity (static ip).
// Each of these labels is a record of intent to do
// something. The controller sets these labels and each
// member watches for them and takes the appropriate
// actions.
//
// See the sidecar design doc for more details.
const (
	// SeedLabel determines if a member is a seed or not.
	SeedLabel = "scylla/seed"

	// IpLabel determines ip of pod when using hostNetwork
	// Used for pod replacement, seed ip and listen/broadcast ip
	IpLabel = "scylla/ip"

	// DecommissionLabel expresses the intent to decommission
	// the specific member. The presence of the label expresses
	// the intent to decommission. If the value is true, it means
	// the member has finished decommissioning.
	// Values: {true, false}
	DecommissionLabel = "scylla/decommissioned"

	// ReplaceLabel express the intent to replace pod under the specific member.
	ReplaceLabel = "scylla/replace"

	// NodeMaintenanceLabel means that node is under maintenance.
	// Readiness check will always fail when this label is added to member service.
	NodeMaintenanceLabel = "scylla/node-maintenance"

	LabelValueTrue  = "true"
	LabelValueFalse = "false"
)

// Generic Labels used on objects created by the operator.
const (
	ClusterNameLabel    = "scylla/cluster"
	DatacenterNameLabel = "scylla/datacenter"
	RackNameLabel       = "scylla/rack"
	MultiDcSeedLabel    = "scylla/multi-dc-seed"

	AppName         = "scylla"
	OperatorAppName = "scylla-operator"
	ManagerAppName  = "scylla-manager"

	PrometheusScrapeAnnotation = "prometheus.io/scrape"
	PrometheusPortAnnotation   = "prometheus.io/port"
)

// Environment Variables
const (
	EnvVarEnvVarPodName = "POD_NAME"
	EnvVarPodNamespace  = "POD_NAMESPACE"
	EnvVarCPU           = "CPU"
	EnvVarRackFromNode  = "RACK_FROM_NODE"
)

// Recorder Values
const (
	// SuccessSynced is used as part of the Event 'reason' when a Cluster is
	// synced.
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a
	// Cluster fails to sync due to a resource of the same name already
	// existing.
	ErrSyncFailed = "ErrSyncFailed"
)

// Bootstrap Values
const (
	// BoostrapOngoing indicate cluster is still bootstraping from multi dc seeds
	BoostrapOngoing = "ongoing"
	// BoostrapFinished indicate cluster has finished bootstraping from multi dc seeds
	BoostrapFinished = "finished"
)

// Configuration Values
const (
	ScyllaContainerName          = "scylla"
	SidecarInjectorContainerName = "sidecar-injection"

	PVCTemplateName = "data"

	SharedDirName = "/mnt/shared"

	ScyllaConfigDirName          = "/mnt/scylla-config"
	ScyllaAgentConfigDirName     = "/mnt/scylla-agent-config"
	ScyllaAgentConfigDefaultFile = "/etc/scylla-manager-agent/scylla-manager-agent.yaml"
	ScyllaClientConfigDirName    = "/mnt/scylla-client-config"
	ScyllaClientConfigFileName   = "scylla-client.yaml"
	ScyllaConfigName             = "scylla.yaml"
	ScyllaRackDCPropertiesName   = "cassandra-rackdc.properties"
	ScyllaIOPropertiesName       = "io_properties.yaml"

	DataDir = "/var/lib/scylla"

	ReadinessProbePath = "/readyz"
	LivenessProbePath  = "/healthz"
	ProbePort          = 8080
	MetricsPort        = 8081
)
