package route

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"code.cloudfoundry.org/routing-api/models"
	"github.com/uber-go/zap"
)

type Endpoint struct {
	sync.RWMutex
	ApplicationId        string
	addr                 string
	Tags                 map[string]string
	ServerCertDomainSAN  string
	PrivateInstanceId    string
	StaleThreshold       time.Duration
	RouteServiceUrl      string
	PrivateInstanceIndex string
	ModificationTag      models.ModificationTag
	Stats                *Stats
	IsolationSegment     string
	useTls               bool
	UpdatedAt            time.Time
}

type EndpointOpts struct {
	AppId                   string
	Host                    string
	Port                    uint16
	ServerCertDomainSAN     string
	PrivateInstanceId       string
	PrivateInstanceIndex    string
	Tags                    map[string]string
	StaleThresholdInSeconds int
	RouteServiceUrl         string
	ModificationTag         models.ModificationTag
	IsolationSegment        string
	UseTLS                  bool
	UpdatedAt               time.Time
}

func NewEndpoint(opts *EndpointOpts) *Endpoint {
	return &Endpoint{
		ApplicationId:        opts.AppId,
		addr:                 fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Tags:                 opts.Tags,
		useTls:               opts.UseTLS,
		ServerCertDomainSAN:  opts.ServerCertDomainSAN,
		PrivateInstanceId:    opts.PrivateInstanceId,
		PrivateInstanceIndex: opts.PrivateInstanceIndex,
		StaleThreshold:       time.Duration(opts.StaleThresholdInSeconds) * time.Second,
		RouteServiceUrl:      opts.RouteServiceUrl,
		ModificationTag:      opts.ModificationTag,
		Stats:                NewStats(),
		IsolationSegment:     opts.IsolationSegment,
		UpdatedAt:            opts.UpdatedAt,
	}
}

func (e *Endpoint) IsTLS() bool {
	return e.useTls
}

func (e *Endpoint) MarshalJSON() ([]byte, error) {
	var jsonObj struct {
		Address             string            `json:"address"`
		TLS                 bool              `json:"tls"`
		TTL                 int               `json:"ttl"`
		RouteServiceUrl     string            `json:"route_service_url,omitempty"`
		Tags                map[string]string `json:"tags"`
		IsolationSegment    string            `json:"isolation_segment,omitempty"`
		PrivateInstanceId   string            `json:"private_instance_id,omitempty"`
		ServerCertDomainSAN string            `json:"server_cert_domain_san,omitempty"`
	}

	jsonObj.Address = e.addr
	jsonObj.TLS = e.IsTLS()
	jsonObj.RouteServiceUrl = e.RouteServiceUrl
	jsonObj.TTL = int(e.StaleThreshold.Seconds())
	jsonObj.Tags = e.Tags
	jsonObj.IsolationSegment = e.IsolationSegment
	jsonObj.PrivateInstanceId = e.PrivateInstanceId
	jsonObj.ServerCertDomainSAN = e.ServerCertDomainSAN
	return json.Marshal(jsonObj)
}

func (e *Endpoint) CanonicalAddr() string {
	return e.addr
}

func (rm *Endpoint) Component() string {
	return rm.Tags["component"]
}

func (e *Endpoint) ToLogData() []zap.Field {
	return []zap.Field{
		zap.String("ApplicationId", e.ApplicationId),
		zap.String("Addr", e.addr),
		zap.Object("Tags", e.Tags),
		zap.String("RouteServiceUrl", e.RouteServiceUrl),
	}
}

func (e *Endpoint) modificationTagSameOrNewer(other *Endpoint) bool {
	return e.ModificationTag == other.ModificationTag || e.ModificationTag.SucceededBy(&other.ModificationTag)
}
