// This file was automatically generated by informer-gen

package v1alpha1

import (
	time "time"

	eventing_kyma_io_v1alpha1 "github.com/kyma-project/kyma/components/event-bus/api/push/eventing.kyma.cx/v1alpha1"
	versioned "github.com/kyma-project/kyma/components/event-bus/generated/push/clientset/versioned"
	internalinterfaces "github.com/kyma-project/kyma/components/event-bus/generated/push/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kyma-project/kyma/components/event-bus/generated/push/listers/eventing.kyma.cx/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// SubscriptionInformer provides access to a shared informer and lister for
// Subscriptions.
type SubscriptionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SubscriptionLister
}

type subscriptionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSubscriptionInformer constructs a new informer for Subscription type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSubscriptionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSubscriptionInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredSubscriptionInformer constructs a new informer for Subscription type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSubscriptionInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EventingV1alpha1().Subscriptions(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EventingV1alpha1().Subscriptions(namespace).Watch(options)
			},
		},
		&eventing_kyma_io_v1alpha1.Subscription{},
		resyncPeriod,
		indexers,
	)
}

func (f *subscriptionInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSubscriptionInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *subscriptionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&eventing_kyma_io_v1alpha1.Subscription{}, f.defaultInformer)
}

func (f *subscriptionInformer) Lister() v1alpha1.SubscriptionLister {
	return v1alpha1.NewSubscriptionLister(f.Informer().GetIndexer())
}
