package ipam

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/hansthienpondt/nipam/pkg/table"
	ipamv1alpha1 "github.com/nokia/k8s-ipam/apis/alloc/ipam/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type Getter interface {
	GetIPAllocation(ctx context.Context) error
}

type GetterConfig struct {
	alloc *ipamv1alpha1.IPAllocation
	rib   *table.RIB
}

func NewGetter(c *GetterConfig) Getter {
	return &getter{
		alloc: c.alloc,
		rib:   c.rib,
	}
}

type getter struct {
	alloc *ipamv1alpha1.IPAllocation
	rib   *table.RIB
	l     logr.Logger
}

func (r *getter) GetIPAllocation(ctx context.Context) error {
	r.l = log.FromContext(ctx).WithValues("name", r.alloc.GetName(), "kind", r.alloc.GetPrefixKind())
	r.l.Info("dynamic allocation")

	labelSelector, err := r.alloc.GetLabelSelector()
	if err != nil {
		return err
	}
	routes := r.rib.GetByLabel(labelSelector)
	if len(routes) != 0 {
		// update the status
		r.alloc.Status.AllocatedPrefix = routes[0].Prefix().String()
		if r.alloc.GetPrefixKind() == ipamv1alpha1.PrefixKindNetwork {
			if !r.alloc.GetCreatePrefix() {
				r.alloc.Status.Gateway = r.getGateway()
			}
		}
	}
	return nil
}

func (r *getter) getGateway() string {
	gatewaySelector, err := r.alloc.GetGatewayLabelSelector()
	if err != nil {
		r.l.Error(err, "cannot get gateway label selector")
		return ""
	}
	r.l.Info("gateway", "gatewaySelector", gatewaySelector)
	routes := r.rib.GetByLabel(gatewaySelector)
	if len(routes) > 0 {
		r.l.Info("gateway", "routes", routes)
		return routes[0].Prefix().Addr().String()
	}
	return ""
}
