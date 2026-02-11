package gmongo

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// Config MongoDB 配置
type Config struct {
	// 基础连接配置
	URI      string `json:"uri" yaml:"uri"`           // MongoDB 连接 URI
	Host     string `json:"host" yaml:"host"`         // 主机地址
	Port     int    `json:"port" yaml:"port"`         // 端口
	Username string `json:"username" yaml:"username"` // 用户名
	Password string `json:"password" yaml:"password"` // 密码
	Database string `json:"database" yaml:"database"` // 数据库名

	// 认证配置
	AuthSource         string `json:"auth_source" yaml:"auth_source"`                   // 认证数据库
	AuthMechanism      string `json:"auth_mechanism" yaml:"auth_mechanism"`             // 认证机制
	ReplicaSet         string `json:"replica_set" yaml:"replica_set"`                   // 副本集名称
	DirectConnection   bool   `json:"direct_connection" yaml:"direct_connection"`       // 直连模式
	RetryWrites        *bool  `json:"retry_writes" yaml:"retry_writes"`                 // 重试写入
	RetryReads         *bool  `json:"retry_reads" yaml:"retry_reads"`                   // 重试读取
	LoadBalanced       bool   `json:"load_balanced" yaml:"load_balanced"`               // 负载均衡
	ServerSelectionTLS bool   `json:"server_selection_tls" yaml:"server_selection_tls"` // TLS 服务器选择

	// 连接池配置
	MaxPoolSize     uint64        `json:"max_pool_size" yaml:"max_pool_size"`         // 最大连接数
	MinPoolSize     uint64        `json:"min_pool_size" yaml:"min_pool_size"`         // 最小连接数
	MaxConnIdleTime time.Duration `json:"max_conn_idle_time" yaml:"max_conn_idle_time"` // 连接最大空闲时间
	MaxConnecting   uint64        `json:"max_connecting" yaml:"max_connecting"`       // 最大连接中数量

	// 超时配置
	ConnectTimeout    time.Duration `json:"connect_timeout" yaml:"connect_timeout"`       // 连接超时
	SocketTimeout     time.Duration `json:"socket_timeout" yaml:"socket_timeout"`         // Socket 超时
	ServerSelectTimeout time.Duration `json:"server_select_timeout" yaml:"server_select_timeout"` // 服务器选择超时
	HeartbeatInterval time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"` // 心跳间隔

	// 读写配置
	ReadPreference  string `json:"read_preference" yaml:"read_preference"`   // 读偏好：primary, primaryPreferred, secondary, secondaryPreferred, nearest
	ReadConcern     string `json:"read_concern" yaml:"read_concern"`         // 读关注：local, available, majority, linearizable, snapshot
	WriteConcern    string `json:"write_concern" yaml:"write_concern"`       // 写关注：majority, w1, w2, w3
	WTimeout        int    `json:"w_timeout" yaml:"w_timeout"`               // 写超时（毫秒）
	Journal         *bool  `json:"journal" yaml:"journal"`                   // 是否写入日志

	// 压缩配置
	Compressors []string `json:"compressors" yaml:"compressors"` // 压缩器：snappy, zlib, zstd

	// 应用配置
	AppName string `json:"app_name" yaml:"app_name"` // 应用名称

	// TLS 配置
	TLS                bool   `json:"tls" yaml:"tls"`                                   // 启用 TLS
	TLSInsecure        bool   `json:"tls_insecure" yaml:"tls_insecure"`                 // 跳过证书验证
	TLSCertificateFile string `json:"tls_certificate_file" yaml:"tls_certificate_file"` // 证书文件路径
	TLSCAFile          string `json:"tls_ca_file" yaml:"tls_ca_file"`                   // CA 文件路径

	// 其他配置
	ZlibLevel       int  `json:"zlib_level" yaml:"zlib_level"`             // Zlib 压缩级别
	ZstdLevel       int  `json:"zstd_level" yaml:"zstd_level"`             // Zstd 压缩级别
	DisableOCSPCheck bool `json:"disable_ocsp_check" yaml:"disable_ocsp_check"` // 禁用 OCSP 检查
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	retryWrites := true
	retryReads := true
	journal := true

	return &Config{
		Host:                "localhost",
		Port:                27017,
		Database:            "test",
		AuthSource:          "admin",
		MaxPoolSize:         100,
		MinPoolSize:         10,
		MaxConnIdleTime:     10 * time.Minute,
		MaxConnecting:       10,
		ConnectTimeout:      10 * time.Second,
		SocketTimeout:       30 * time.Second,
		ServerSelectTimeout: 30 * time.Second,
		HeartbeatInterval:   10 * time.Second,
		ReadPreference:      "primary",
		ReadConcern:         "local",
		WriteConcern:        "majority",
		WTimeout:            5000,
		RetryWrites:         &retryWrites,
		RetryReads:          &retryReads,
		Journal:             &journal,
		Compressors:         []string{"snappy", "zlib", "zstd"},
		ZlibLevel:           6,
		ZstdLevel:           6,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.URI == "" {
		if c.Host == "" {
			return fmt.Errorf("host is required when URI is not provided")
		}
		if c.Port <= 0 || c.Port > 65535 {
			return fmt.Errorf("invalid port: %d", c.Port)
		}
	}

	if c.MaxPoolSize > 0 && c.MinPoolSize > c.MaxPoolSize {
		return fmt.Errorf("min_pool_size (%d) cannot be greater than max_pool_size (%d)", c.MinPoolSize, c.MaxPoolSize)
	}

	return nil
}

// GetURI 获取连接 URI
func (c *Config) GetURI() string {
	if c.URI != "" {
		return c.URI
	}

	// 构建 URI
	uri := "mongodb://"

	// 添加认证信息
	if c.Username != "" {
		uri += c.Username
		if c.Password != "" {
			uri += ":" + c.Password
		}
		uri += "@"
	}

	// 添加主机和端口
	uri += fmt.Sprintf("%s:%d", c.Host, c.Port)

	// 添加数据库
	if c.Database != "" {
		uri += "/" + c.Database
	}

	// 添加查询参数
	params := ""
	if c.AuthSource != "" {
		params += "&authSource=" + c.AuthSource
	}
	if c.AuthMechanism != "" {
		params += "&authMechanism=" + c.AuthMechanism
	}
	if c.ReplicaSet != "" {
		params += "&replicaSet=" + c.ReplicaSet
	}
	if c.DirectConnection {
		params += "&directConnection=true"
	}
	if c.RetryWrites != nil && *c.RetryWrites {
		params += "&retryWrites=true"
	}
	if c.RetryReads != nil && *c.RetryReads {
		params += "&retryReads=true"
	}
	if c.LoadBalanced {
		params += "&loadBalanced=true"
	}
	if c.TLS {
		params += "&tls=true"
	}
	if c.TLSInsecure {
		params += "&tlsInsecure=true"
	}

	if params != "" {
		uri += "?" + params[1:] // 去掉第一个 &
	}

	return uri
}

// ToClientOptions 转换为 mongo.ClientOptions
func (c *Config) ToClientOptions() *options.ClientOptions {
	opts := options.Client().ApplyURI(c.GetURI())

	// 连接池配置
	if c.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(c.MaxPoolSize)
	}
	if c.MinPoolSize > 0 {
		opts.SetMinPoolSize(c.MinPoolSize)
	}
	if c.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(c.MaxConnIdleTime)
	}
	if c.MaxConnecting > 0 {
		opts.SetMaxConnecting(c.MaxConnecting)
	}

	// 超时配置
	if c.ConnectTimeout > 0 {
		opts.SetConnectTimeout(c.ConnectTimeout)
	}
	if c.SocketTimeout > 0 {
		opts.SetSocketTimeout(c.SocketTimeout)
	}
	if c.ServerSelectTimeout > 0 {
		opts.SetServerSelectionTimeout(c.ServerSelectTimeout)
	}
	if c.HeartbeatInterval > 0 {
		opts.SetHeartbeatInterval(c.HeartbeatInterval)
	}

	// 读偏好
	if c.ReadPreference != "" {
		if rp := parseReadPreference(c.ReadPreference); rp != nil {
			opts.SetReadPreference(rp)
		}
	}

	// 读关注
	if c.ReadConcern != "" {
		if rc := parseReadConcern(c.ReadConcern); rc != nil {
			opts.SetReadConcern(rc)
		}
	}

	// 写关注
	if c.WriteConcern != "" {
		if wc := parseWriteConcern(c.WriteConcern, c.WTimeout, c.Journal); wc != nil {
			opts.SetWriteConcern(wc)
		}
	}

	// 压缩配置
	if len(c.Compressors) > 0 {
		opts.SetCompressors(c.Compressors)
	}
	if c.ZlibLevel > 0 {
		opts.SetZlibLevel(c.ZlibLevel)
	}
	if c.ZstdLevel > 0 {
		opts.SetZstdLevel(c.ZstdLevel)
	}

	// 应用名称
	if c.AppName != "" {
		opts.SetAppName(c.AppName)
	}

	// 重试配置
	if c.RetryWrites != nil {
		opts.SetRetryWrites(*c.RetryWrites)
	}
	if c.RetryReads != nil {
		opts.SetRetryReads(*c.RetryReads)
	}

	// 其他配置
	if c.DisableOCSPCheck {
		opts.SetDisableOCSPEndpointCheck(true)
	}

	return opts
}

// parseReadPreference 解析读偏好
func parseReadPreference(rp string) *readpref.ReadPref {
	switch rp {
	case "primary":
		return readpref.Primary()
	case "primaryPreferred":
		return readpref.PrimaryPreferred()
	case "secondary":
		return readpref.Secondary()
	case "secondaryPreferred":
		return readpref.SecondaryPreferred()
	case "nearest":
		return readpref.Nearest()
	default:
		return nil
	}
}

// parseReadConcern 解析读关注
func parseReadConcern(rc string) *readconcern.ReadConcern {
	switch rc {
	case "local":
		return readconcern.Local()
	case "available":
		return readconcern.Available()
	case "majority":
		return readconcern.Majority()
	case "linearizable":
		return readconcern.Linearizable()
	case "snapshot":
		return readconcern.Snapshot()
	default:
		return nil
	}
}

// parseWriteConcern 解析写关注
func parseWriteConcern(wc string, timeout int, journal *bool) *writeconcern.WriteConcern {
	var w writeconcern.WriteConcern

	switch wc {
	case "majority":
		w = *writeconcern.Majority()
	case "w1":
		w = *writeconcern.W1()
	case "w2":
		w = *writeconcern.Unacknowledged() // 需要自定义
		w.W = 2
	case "w3":
		w = *writeconcern.Unacknowledged()
		w.W = 3
	default:
		return nil
	}

	// 设置超时
	if timeout > 0 {
		w.WTimeout = time.Duration(timeout) * time.Millisecond
	}

	// 设置日志
	if journal != nil {
		w.Journal = journal
	}

	return &w
}
