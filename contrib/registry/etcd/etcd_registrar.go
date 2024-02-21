package etcd

import (
	"context"
	gberror "ghostbb.io/gb/errors/gb_error"
	gbsvc "ghostbb.io/gb/net/gb_svc"

	etcd3 "go.etcd.io/etcd/client/v3"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (r *Registry) Register(ctx context.Context, service gbsvc.Service) (gbsvc.Service, error) {
	service = NewService(service)
	r.lease = etcd3.NewLease(r.client)
	grant, err := r.lease.Grant(ctx, int64(r.keepaliveTTL.Seconds()))
	if err != nil {
		return nil, gberror.Wrapf(err, `etcd grant failed with keepalive ttl "%s"`, r.keepaliveTTL)
	}
	var (
		key   = service.GetKey()
		value = service.GetValue()
	)
	_, err = r.client.Put(ctx, key, value, etcd3.WithLease(grant.ID))
	if err != nil {
		return nil, gberror.Wrapf(
			err,
			`etcd put failed with key "%s", value "%s", lease "%d"`,
			key, value, grant.ID,
		)
	}
	r.logger.Debugf(
		ctx,
		`etcd put success with key "%s", value "%s", lease "%d"`,
		key, value, grant.ID,
	)
	keepAliceCh, err := r.client.KeepAlive(context.Background(), grant.ID)
	if err != nil {
		return nil, err
	}
	go r.doKeepAlive(grant.ID, keepAliceCh)
	return service, nil
}

// Deregister off-lines and removes `service` from the Registry.
func (r *Registry) Deregister(ctx context.Context, service gbsvc.Service) error {
	_, err := r.client.Delete(ctx, service.GetKey())
	if r.lease != nil {
		_ = r.lease.Close()
	}
	return err
}

// doKeepAlive continuously keeps alive the lease from ETCD.
func (r *Registry) doKeepAlive(leaseID etcd3.LeaseID, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse) {
	var ctx = context.Background()
	for {
		select {
		case <-r.client.Ctx().Done():
			r.logger.Noticef(ctx, "keepalive done for lease id: %d", leaseID)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				// r.logger.Debugf(ctx, `keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				r.logger.Noticef(ctx, `keepalive exit, lease id: %d`, leaseID)
				return
			}
		}
	}
}
