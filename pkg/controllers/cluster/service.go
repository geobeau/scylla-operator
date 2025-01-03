package cluster

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/v1"
	"github.com/scylladb/scylla-operator/pkg/controllers/cluster/resource"
	"github.com/scylladb/scylla-operator/pkg/naming"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// syncClusterHeadlessService checks if a Headless Service exists
// for the given Cluster, in order for the StatefulSets to utilize it.
// If it doesn't exists, then create it.
func (cc *ClusterReconciler) syncClusterHeadlessService(ctx context.Context, c *scyllav1.ScyllaCluster) error {
	clusterHeadlessService := resource.HeadlessServiceForCluster(c)
	_, err := controllerutil.CreateOrUpdate(ctx, cc.Client, clusterHeadlessService, serviceMutateFn(ctx, clusterHeadlessService, cc.Client))
	if err != nil {
		return errors.Wrapf(err, "error syncing headless service %s", clusterHeadlessService.Name)
	}
	return nil
}

// syncMemberServices checks, for every Pod of the Cluster that
// has been created, if a corresponding ClusterIP Service exists,
// which will serve as a static ip.
// If it doesn't exist, it creates it.
// It also assigns the first two members of each rack as seeds.
func (cc *ClusterReconciler) syncMemberServices(ctx context.Context, c *scyllav1.ScyllaCluster) error {
	podlist := &corev1.PodList{}

	// For every Pod of the cluster that exists, check that a
	// a corresponding ClusterIP Service exists, and if it doesn't,
	// create it.
	for _, r := range c.Spec.Datacenter.Racks {
		// Get all Pods for this rack
		opts := client.MatchingLabels(naming.RackLabels(r, c))
		err := cc.List(ctx, podlist, opts)
		if err != nil {
			return errors.Wrapf(err, "listing pods for rack %s failed", r.Name)
		}

		memberCount := r.Members
		maxIndex := memberCount - 1

		for _, pod := range podlist.Items {
			memberService := resource.MemberServiceForPod(&pod, c)
			svcIndex, err := naming.IndexFromName(memberService.Name)
			if err != nil {
				return errors.WithStack(err)
			}

			if svcIndex <= maxIndex {
				op, err := controllerutil.CreateOrUpdate(ctx, cc.Client, memberService, serviceMemberMutateFn(memberService, memberService.DeepCopy()))
				if err != nil {
					return errors.Wrapf(err, "error syncing member service %s", memberService.Name)
				}
				switch op {
				case controllerutil.OperationResultCreated:
					cc.Logger.Info(ctx, "Member service created", "member", memberService.Name, "labels", memberService.Labels)
				case controllerutil.OperationResultUpdated:
					cc.Logger.Info(ctx, "Member service updated", "member", memberService.Name, "labels", memberService.Labels)
				}
			} else {
				cc.Logger.Debug(ctx, "Member service not created as index greater than max members", "svcIndex", svcIndex, "rackName", r.Name, "podName", pod.Name)
			}
		}
	}
	return nil
}

func serviceMemberMutateFn(service *corev1.Service, newService *corev1.Service) func() error {
	return func() error {
		// Ensure scylla/ip label is well updated
		if ip, ok := newService.ObjectMeta.Labels[naming.IpLabel]; ok && ip != "" {
			service.ObjectMeta.Labels[naming.IpLabel] = ip
		}

		return nil
	}
}

// syncService checks if the given Service exists and creates it if it doesn't
// it creates it
func serviceMutateFn(ctx context.Context, newService *corev1.Service, client client.Client) func() error {
	return func() error {
		// TODO: probably nothing has to be done, check v1 implementation of CreateOrUpdate
		//existingService := existing.(*corev1.Service)
		//if !reflect.DeepEqual(newService.Spec, existingService.Spec) {
		//	return client.Update(ctx, existing)
		//}
		return nil
	}
}

func (cc *ClusterReconciler) syncMultiDcServices(ctx context.Context, cluster *scyllav1.ScyllaCluster) error {
	for id, seed := range cluster.Spec.MultiDcCluster.Seeds {
		multiDcServiceName := fmt.Sprintf("%s-%s-multi-dc-seed-%d", cluster.Name, cluster.Spec.Datacenter.Name, id)

		cc.Logger.Info(ctx, "Create multi dc seed", "multiDcServiceName", multiDcServiceName)
		multiDcService := resource.ServiceForMultiDcSeed(multiDcServiceName, seed, cluster)
		op, err := controllerutil.CreateOrUpdate(ctx, cc.Client, multiDcService, serviceMultiDcMutateFn(multiDcService, multiDcService.DeepCopy()))
		if err != nil {
			return errors.Wrapf(err, "error syncing multi dc service %s", multiDcService.Name)
		}
		switch op {
		case controllerutil.OperationResultCreated:
			cc.Logger.Info(ctx, "Multi Dc seed service created", "multiDcSeed", multiDcService.Name, "labels", multiDcService.Labels)
		case controllerutil.OperationResultUpdated:
			cc.Logger.Info(ctx, "Multi Dc seed service updated", "multiDcSeed", multiDcService.Name, "labels", multiDcService.Labels)
		}
	}
	return nil
}

func serviceMultiDcMutateFn(service *corev1.Service, newService *corev1.Service) func() error {
	return func() error {
		service.ObjectMeta.Labels[naming.IpLabel] = newService.ObjectMeta.Labels[naming.IpLabel]
		return nil
	}
}
