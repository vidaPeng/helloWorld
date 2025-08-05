package config

type CloudResource struct {
	CloudName string    `json:"cloud_name" form:"cloud_name"` // 云厂商名字
	CloudType string    `json:"cloud_type" form:"cloud_type"` // 云资源的类型
	CloudData CloudData `json:"cloud_data" form:"cloud_data"` // 申请的云资源所需要的字段
}

type CloudData struct {
	// 前三个是标准的需要的字段
	Name    string `json:"name"`
	ENV     string `json:"env"`
	PixccId string `json:"pixcc_id"`
	// 后面是各个调用所需要配置的参数
	MySqlData MySqlData `json:"mysql_data,omitempty"`
	RedisData RedisData `json:"redis_data,omitempty"`
	VmData    VmData    `json:"vm_data,omitempty"`
	S3Data    S3Data    `json:"s3_data,omitempty"`
	//DnsData        DnsData        `json:"dns_data,omitempty"`
}

type MySqlData struct {
	Specification        string `json:"specification"`  // 数据库的规格
	EnableCluster        bool   `json:"enable_cluster"` // 是否启用集群
	AdminUser            string `json:"admin_user"`
	AdminPassword        string `json:"admin_password"`
	DataStorageSizeInGBs int    `json:"data_storage_size_in_gbs"`
	MysqlVersion         string `json:"mysql_version"`
	ProjectName          string `json:"project_name"`

	ProjectId string `json:"project_id"` // 项目 ID
}

type RedisData struct {
	EnableSharding  bool   `json:"enable_sharding"`  // 是否开启分片
	NodeMemory      int    `json:"node_memory"`      // 每个节点的内存大小 G
	SoftwareVersion string `json:"software_version"` // redis 版本
	NodeCount       int    `json:"node_count"`       // 节点数量
	ShardCount      int    `json:"shard_count"`      // 分片数量
	ProjectName     string `json:"project_name"`

	ProjectId string `json:"project_id"` // 项目 ID
}

type VmData struct {
	ProcessorType  string `json:"processor_type"`  // 处理器类型
	SystemType     string `json:"system_type"`     // 系统类型 Linux or win
	Specification  string `json:"specification"`   // 规格大小
	AdditionalDisk int    `json:"additional_disk"` // 附加磁盘的大小，单位转为 G
	ProjectName    string `json:"project_name"`

	ProjectId string `json:"project_id"` // 项目 ID
}

type S3Data struct {
	ProjectId string `json:"project_id"`
}

type Response struct {
	Data    ProjectResp `json:"data"`
	Message string      `json:"message"`
}

type ProjectResp struct {
	Count    int       `json:"count"`
	Projects []Project `json:"projects"`
}

type Project struct {
	ProjectId string `json:"project_id"`
	Name      string `json:"name"`
	State     string `json:"state"`
}

// APIResponse 是最外层的结构，对应整个 JSON
type APIResponse struct {
	// "data" 字段是一个 map，键是数据库实例名 (string)，值是实例的详细信息 (DBInstance)
	Data    map[string]BucketInfo `json:"data"`
	Message string                `json:"message"`
}

// DBInstance 包含了每个数据库实例的具体属性
type DBInstance struct {
	BinlogEnabled   bool     `json:"binlogEnabled"`
	CharacterSet    string   `json:"characterSet"`
	ClusterRole     string   `json:"clusterRole"`
	ConnectionName  string   `json:"connectionName"`
	CreateTime      string   `json:"createTime"`
	DatabaseVersion string   `json:"databaseVersion"`
	DBName          string   `json:"dbName"`
	IPAddresses     []string `json:"ipAddresses"`
	MaxConnections  string   `json:"maxConnections"` // 注意：这个是 string 类型
	Operator        string   `json:"operator"`
	Port            int      `json:"port"`
	Region          string   `json:"region"`
	Status          string   `json:"status"`
	Tier            string   `json:"tier"`

	ProjectID string `json:"projectID"` // 添加 Project ID 信息
}

type RedisInfo struct {
	CreateTime   string `json:"createTime"`
	DisplayName  string `json:"displayName"` // 用户在控制台看到的名称
	Host         string `json:"host"`
	InstanceName string `json:"instanceName"`
	Location     string `json:"location"` // 地域
	MemorySizeGb int    `json:"memorySizeGb"`
	Operator     string `json:"operator"`
	Port         int    `json:"port"`
	RedisVersion string `json:"redisVersion"`
	Status       string `json:"status"`
	Tier         string `json:"tier"` // 实例的服务层级

	ProjectID string `json:"projectID"` // 添加 Project ID 信息
}

type VmInfo struct {
	CreationTimestamp string            `json:"creationTimestamp"`
	ExternalIP        string            `json:"externalIp"`
	InstanceID        uint64            `json:"instanceId"` // 使用,string确保大整数能被正确解析
	InstanceName      string            `json:"instanceName"`
	InternalIP        string            `json:"internalIp"`
	Labels            map[string]string `json:"labels"`
	MachineType       string            `json:"machineType"`
	Status            string            `json:"status"`
	Zone              string            `json:"zone"`

	ProjectID string `json:"projectID"` // 添加 Project ID 信息
}

type BucketInfo struct {
	BucketName        string            `json:"bucketName"`
	CreateTime        string            `json:"createTime"`
	Labels            map[string]string `json:"labels"`
	Location          string            `json:"location"`
	StorageClass      string            `json:"storageClass"`
	VersioningEnabled bool              `json:"versioningEnabled"`

	ProjectID string `json:"projectID"` // 添加 Project ID 信息
}
