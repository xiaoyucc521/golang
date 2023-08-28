package weight

import (
	"math/rand"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

// Name is the name of round_robin balancer.
const Name = "weight"

// AttributeKey is the type used as the key to store AddrInfo in the Attributes
// field of resolver.Address.
type attributeKey struct{}

// AddrInfo will be stored inside Address metadata in order to use weighted balancer.
type AddrInfo struct {
	Weight int
}

// SetAddrInfo returns a copy of addr in which the Attributes field is updated
// with addrInfo.
func SetAddrInfo(addr resolver.Address, addrInfo AddrInfo) resolver.Address {
	addr.Attributes = attributes.New(attributeKey{}, addrInfo)
	return addr
}

// GetAddrInfo returns the AddrInfo stored in the Attributes fields of addr.
func GetAddrInfo(addr resolver.Address) AddrInfo {
	v := addr.Attributes.Value(attributeKey{})
	ai, _ := v.(AddrInfo)
	return ai
}

// newBuilder creates a new roundrobin balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &rrPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type rrPickerBuilder struct{}

func (*rrPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc, addr := range info.ReadySCs {

		node := GetAddrInfo(addr.Address)
		// 简单的负载均衡算法
		if node.Weight >= 0 {
			for i := 0; i < node.Weight; i++ {
				scs = append(scs, sc)
			}
		} else {
			scs = append(scs, sc)
		}

	}
	return &rrPicker{
		subConns: scs,
	}
}

type rrPicker struct {
	subConns []balancer.SubConn
}

func (p *rrPicker) Pick(balancer.PickInfo) (balancer.PickResult, error) {
	subConnsLen := len(p.subConns)
	nextIndex := rand.Intn(subConnsLen)
	sc := p.subConns[nextIndex]
	return balancer.PickResult{SubConn: sc}, nil
}
