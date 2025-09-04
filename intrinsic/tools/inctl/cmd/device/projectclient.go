// Copyright 2023 Intrinsic Innovation LLC

package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"google.golang.org/grpc"
	"intrinsic/frontend/cloud/devicemanager/shared/shared"
	"intrinsic/tools/inctl/auth/auth"

	clustermanagergrpcpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
	clustermanagerpb "intrinsic/frontend/cloud/api/v1/clustermanager_api_go_grpc_proto"
)

var (
	// These will be returned on corresponding http error codes, since they are errors that are
	// expected and can be printed with better UX than just the number.
	errNotFound     = fmt.Errorf("Not found")
	errBadGateway   = fmt.Errorf("Bad Gateway")
	errUnauthorized = fmt.Errorf("Unauthorized")
)

// authedClient injects an api key for the project into every request.
type authedClient struct {
	client       *http.Client
	baseURL      url.URL
	projectName  string
	organization string
	grpcConn     *grpc.ClientConn
	grpcClient   clustermanagergrpcpb.ClustersServiceClient
}

// newClient returns a http.Client compatible that injects auth for the project into every request.
func newClient(ctx context.Context, projectName string, orgName string, clusterName string) (authedClient, error) {
	// create a cloud connection to the cluster via the relay with a callback to get the token source
	opts := []auth.ConnectionOptsFunc{
		auth.WithProject(projectName), auth.WithOrg(orgName), auth.WithCluster(clusterName),
	}
	conn, err := auth.NewCloudConnection(ctx, opts...)
	if err != nil {
		return authedClient{}, err
	}
	// create a http client from the cloud connection
	cl, err := auth.NewCloudClient(ctx, opts...)
	if err != nil {
		return authedClient{}, err
	}

	return authedClient{
		client: cl,
		baseURL: url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("www.endpoints.%s.cloud.goog", projectName),
			Path:   "/api/devices/",
		},
		projectName:  projectName,
		organization: orgName,
		grpcConn:     conn,
		grpcClient:   clustermanagergrpcpb.NewClustersServiceClient(conn),
	}, nil
}

// close closes the grpc connection if it exists.
func (c *authedClient) close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

func (c *authedClient) getStatusNetwork(ctx context.Context, clusterName, deviceID string) (map[string]shared.StatusInterface, error) {
	req := &clustermanagerpb.GetStatusRequest{
		Project:   c.projectName,
		Org:       c.organization,
		ClusterId: clusterName,
		DeviceId:  deviceID,
	}
	resp, err := c.grpcClient.GetStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	statusNetwork := map[string]shared.StatusInterface{}
	for in, ifa := range resp.GetInterfaces() {
		statusNetwork[in] = shared.StatusInterface{
			IPAddress: ifa.GetAddresses(),
		}
	}
	return statusNetwork, nil
}

func translateNetworkConfig(n *clustermanagerpb.IntOSNetworkConfig) map[string]shared.Interface {
	configMap := map[string]shared.Interface{}
	for name, inf := range n.GetInterfaces() {
		ns := inf.GetNameservers()
		configMap[name] = shared.Interface{
			DHCP4:    inf.GetDhcp4(),
			Gateway4: inf.GetGateway4(),
			DHCP6:    &inf.Dhcp6,
			Gateway6: inf.GetGateway6(),
			MTU:      int64(inf.GetMtu()),
			Nameservers: shared.Nameservers{
				Search:    ns.GetSearch(),
				Addresses: ns.GetAddresses(),
			},
			Addresses: inf.GetAddresses(),
			Realtime:  inf.GetRealtime(),
			EtherType: int64(inf.GetEtherType()),
		}
	}
	return configMap
}

func translateToNetworkConfig(n map[string]shared.Interface) *clustermanagerpb.IntOSNetworkConfig {
	c := &clustermanagerpb.IntOSNetworkConfig{
		Interfaces: make(map[string]*clustermanagerpb.IntOSInterfaceConfig),
	}
	for name, inf := range n {
		dhcp6 := false
		if inf.DHCP6 != nil {
			dhcp6 = *inf.DHCP6
		}
		conf := &clustermanagerpb.IntOSInterfaceConfig{
			Dhcp4:    inf.DHCP4,
			Gateway4: inf.Gateway4,
			Dhcp6:    dhcp6,
			Gateway6: inf.Gateway6,
			Mtu:      int32(inf.MTU),
			Nameservers: &clustermanagerpb.NameserverConfig{
				Search:    inf.Nameservers.Search,
				Addresses: inf.Nameservers.Addresses,
			},
			Addresses: inf.Addresses,
			Realtime:  inf.Realtime,
		}
		switch inf.EtherType {
		default:
			conf.EtherType = clustermanagerpb.IntOSInterfaceConfig_ETHER_TYPE_UNSPECIFIED
		case shared.EtherTypeEtherCAT:
			conf.EtherType = clustermanagerpb.IntOSInterfaceConfig_ETHER_TYPE_ETHERCAT
		}
		c.Interfaces[name] = conf
	}
	return c
}

func (c *authedClient) getNetworkConfig(ctx context.Context, clusterName, deviceID string) (map[string]shared.Interface, error) {
	req := &clustermanagerpb.GetNetworkConfigRequest{
		Project: c.projectName,
		Org:     c.organization,
		Cluster: clusterName,
		Device:  deviceID,
	}
	resp, err := c.grpcClient.GetNetworkConfig(ctx, req)
	if err != nil {
		return nil, err
	}
	return translateNetworkConfig(resp), nil
}

// postDevice acts similar to [http.Post] but takes a context and injects base path of the device manager for the project.
func (c *authedClient) postDevice(ctx context.Context, cluster, deviceID, subPath string, body io.Reader) (*http.Response, error) {
	reqURL := c.baseURL

	reqURL.Path = filepath.Join(reqURL.Path, subPath)
	reqURL.RawQuery = url.Values{"device-id": []string{deviceID}, "cluster": []string{cluster}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

// getDevice acts similar to [http.Get] but takes a context and injects base path of the device manager for the project.
func (c *authedClient) getDevice(ctx context.Context, cluster, deviceID, subPath string) (*http.Response, error) {
	reqURL := c.baseURL

	reqURL.Path = filepath.Join(reqURL.Path, subPath)
	reqURL.RawQuery = url.Values{"device-id": []string{deviceID}, "cluster": []string{cluster}}.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

// getJSON acts similar to [GetDevice] but also does [json.Decode] and enforces [http.StatusOK].
func (c *authedClient) getJSON(ctx context.Context, cluster, deviceID, subPath string, value any) error {
	resp, err := c.getDevice(ctx, cluster, deviceID, subPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return errNotFound
		}
		if resp.StatusCode == http.StatusBadGateway {
			return errBadGateway
		}
		if resp.StatusCode == http.StatusUnauthorized {
			return errUnauthorized
		}

		return fmt.Errorf("get status code: %v", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(value)
}
