package model

type KubeList struct {
	BaseList
	Items []*Kube `json:"items"`
}

// Kube objects contains global info about kubernetes ckusters.
type Kube struct {
	BaseModel

	// belongs_to CloudAccount
	CloudAccount     *CloudAccount `json:"cloud_account,omitempty" gorm:"ForeignKey:CloudAccountName;AssociationForeignKey:Name"`
	CloudAccountName string        `json:"cloud_account_name" validate:"nonzero" gorm:"not null;index" sg:"immutable"`

	// has_many Nodes
	Nodes []*Node `json:"nodes,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`

	// has_many LoadBalancers
	LoadBalancers []*LoadBalancer `json:"load_balancers,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`

	// has_many KubeResources
	KubeResources []*KubeResource `json:"kube_resources,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`

	// has_many HelmReleases
	HelmReleases []*HelmRelease `json:"helm_releases,omitempty" gorm:"ForeignKey:KubeName;AssociationForeignKey:Name"`

	Name string `json:"name" validate:"nonzero,max=12,regexp=^[a-z]([-a-z0-9]*[a-z0-9])?$" gorm:"not null;unique_index" sg:"immutable"`

	MasterNodeSize string `json:"master_node_size" validate:"nonzero" sg:"immutable"`

	NodeSizes     []string `json:"node_sizes" gorm:"-" validate:"min=1" sg:"store_as_json_in=NodeSizesJSON"`
	NodeSizesJSON []byte   `json:"-" gorm:"not null"`

	Username string `json:"username" validate:"nonzero" sg:"immutable"`
	Password string `json:"password" validate:"nonzero" sg:"immutable"`

	HeapsterVersion          string `json:"heapster_version" validate:"nonzero" sg:"default=v1.1.0,immutable"`
	HeapsterMetricResolution string `json:"heapster_metric_resolution" validate:"regexp=^([0-9]+[smhd])+$" sg:"default=20s,immutable"`

	// NOTE due to how we marshal this as JSON, it's difficult to have this stored
	// as an interface, because unmarshalling causes us to lose the underlying
	// type. So, this is kindof like a whacky form of single-table inheritance.
	AWSConfig     *AWSKubeConfig `json:"aws_config,omitempty" gorm:"-" sg:"store_as_json_in=AWSConfigJSON,immutable"`
	AWSConfigJSON []byte         `json:"-"`

	DigitalOceanConfig     *DOKubeConfig `json:"digitalocean_config,omitempty" gorm:"-" sg:"store_as_json_in=DigitalOceanConfigJSON,immutable"`
	DigitalOceanConfigJSON []byte        `json:"-"`

	OpenStackConfig     *OSKubeConfig `json:"openstack_config,omitempty" gorm:"-" sg:"store_as_json_in=OpenStackConfigJSON,immutable"`
	OpenStackConfigJSON []byte        `json:"-"`

	GCEConfig     *GCEKubeConfig `json:"gce_config,omitempty" gorm:"-" sg:"store_as_json_in=GCEConfigJSON,immutable"`
	GCEConfigJSON []byte         `json:"-"`

	MasterPublicIP string `json:"master_public_ip" sg:"readonly"`

	Ready bool `json:"ready" sg:"readonly" gorm:"index"`
}

// AWSKubeConfig holds aws specific information about AWS based KUbernetes clusters.
type AWSKubeConfig struct {
	Region           string `json:"region" validate:"nonzero,regexp=^[a-z]{2}-[a-z]+-[0-9]$"`
	AvailabilityZone string `json:"availability_zone"`
	VPCIPRange       string `json:"vpc_ip_range" validate:"nonzero" sg:"default=172.20.0.0/16"`
	// TODO this should be a slice of objects instead of maps, since we have a rigid key structure
	PublicSubnetIPRange []map[string]string `json:"public_subnet_ip_range"`
	KubeMasterCount     int                 `json:"kube_master_count"`
	MultiAZ             bool                `json:"multi_az"`
	SSHPubKey           string              `json:"ssh_pub_key"`
	BucketName          string              `json:"bucket_name,omitempty" sg:"readonly"`
	NodeVolumeSize      int                 `json:"node_volume_size" sg:"default=100"`
	MasterVolumeSize    int                 `json:"master_volume_size" sg:"default=100"`
	KubernetesVersion   string              `json:"kubernetes_version" validate:"nonzero" sg:"default=1.5.1"`

	MasterPrivateIP               string   `json:"master_private_ip" sg:"readonly"`
	LastSelectedAZ                string   `json:"last_selected_az" sg:"readonly"` // if using multiAZ this is the last az the node build used.
	ETCDDiscoveryURL              string   `json:"etcd_discovery_url" sg:"readonly"`
	PrivateKey                    string   `json:"private_key,omitempty" sg:"readonly"`
	VPCID                         string   `json:"vpc_id"`
	VPCMANAGED                    bool     `json:"vpc_managed"`
	InternetGatewayID             string   `json:"internet_gateway_id" sg:"readonly"`
	RouteTableID                  string   `json:"route_table_id" sg:"readonly"`
	RouteTableSubnetAssociationID []string `json:"route_table_subnet_association_id" sg:"readonly"`
	ELBSecurityGroupID            string   `json:"elb_security_group_id" sg:"readonly"`
	NodeSecurityGroupID           string   `json:"node_security_group_id" sg:"readonly"`
	MasterID                      string   `json:"master_id" sg:"readonly"`
	MasterNodes                   []string `json:"master_nodes" sg:"readonly"`
}

// DOKubeConfig holds do specific information about DO based KUbernetes clusters.
type DOKubeConfig struct {
	Region            string `json:"region" validate:"nonzero"`
	SSHKeyFingerprint string `json:"ssh_key_fingerprint" validate:"nonzero"`

	MasterID int `json:"master_id" sg:"readonly"`
}

// OSKubeConfig holds do specific information about Open Stack based KUbernetes clusters.
type OSKubeConfig struct {
	Region             string `json:"region" validate:"nonzero"`
	SSHPubKey          string `json:"ssh_pub_key" validate:"nonzero"`
	PrivateSubnetRange string `json:"private_subnet_ip_range" validate:"nonzero" sg:"default=172.20.0.0/24"`
	PublicGatwayID     string `json:"public_gateway_id" validate:"nonzero" sg:"default=disabled"`

	MasterID        string `json:"master_id" sg:"readonly"`
	MasterPrivateIP string `json:"master_private_ip" sg:"readonly"`
	NetworkID       string `json:"network_id" sg:"readonly"`
	SubnetID        string `json:"subnet_id" sg:"readonly"`
	RouterID        string `json:"router_id" sg:"readonly"`
	FloatingIpID    string `json:"floating_ip_id" sg:"readonly"`
}

// GCEKubeConfig holds do specific information about DO based KUbernetes clusters.
type GCEKubeConfig struct {
	Zone                string   `json:"zone" validate:"nonzero"`
	MasterInstanceGroup string   `json:"instance_group" sg:"readonly"`
	MinionInstanceGroup string   `json:"instance_group" sg:"readonly"`
	MasterNodes         []string `json:"master_nodes" sg:"readonly"`
	MasterName          string   `json:"master_name" sg:"readonly"`
	KubeMasterCount     int      `json:"kube_master_count"`

	// Template vars
	SSHPubKey         string `json:"ssh_pub_key" validate:"nonzero"`
	KubernetesVersion string `json:"kubernetes_version" validate:"nonzero" sg:"default=1.5.1"`
	ETCDDiscoveryURL  string `json:"etcd_discovery_url" sg:"readonly"`
	MasterPrivateIP   string `json:"master_private_ip" sg:"readonly"`
}
