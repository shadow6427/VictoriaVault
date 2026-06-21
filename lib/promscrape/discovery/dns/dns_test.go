package dns

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/VictoriaMetrics/VictoriaMetrics/lib/netutil"
)

type failingResolver struct{}

func (failingResolver) LookupSRV(context.Context, string, string, string) (string, []*net.SRV, error) {
	return "", nil, errors.New("lookup failed")
}

func (failingResolver) LookupIPAddr(context.Context, string) ([]net.IPAddr, error) {
	return nil, errors.New("lookup failed")
}

func (failingResolver) LookupMX(context.Context, string) ([]*net.MX, error) {
	return nil, errors.New("lookup failed")
}

func TestGetLabelsReturnsErrorIfAllLookupsFail(t *testing.T) {
	resolverOld := netutil.Resolver
	netutil.Resolver = failingResolver{}
	defer func() {
		netutil.Resolver = resolverOld
	}()

	port := 9100
	f := func(sdc *SDConfig) {
		t.Helper()
		if labels, err := sdc.GetLabels(""); err == nil {
			t.Fatalf("expecting non-nil error for config %#v; got labels=%v", sdc, labels)
		}
	}

	f(&SDConfig{
		Names: []string{"example.com", "example.net"},
		Type:  "SRV",
	})
	f(&SDConfig{
		Names: []string{"example.com", "example.net"},
		Type:  "MX",
	})
	f(&SDConfig{
		Names: []string{"example.com", "example.net"},
		Type:  "A",
		Port:  &port,
	})
}
