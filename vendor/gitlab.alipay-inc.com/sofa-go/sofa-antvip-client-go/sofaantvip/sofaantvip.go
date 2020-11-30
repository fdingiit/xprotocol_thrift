package sofaantvip

type Listener interface {
	OnDomainChanged(err error, name string, domain *VipDomain)
}

type ListenerFunc func(err error, name string, domain *VipDomain)

func (l ListenerFunc) OnDomainChanged(err error, name string, domain *VipDomain) {
	l(err, name, domain)
}
