package fake

import (
	v1alpha1 "github.com/kyma-project/kyma/components/event-bus/api/push/eventing.kyma.cx/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSubscriptions implements SubscriptionInterface
type FakeSubscriptions struct {
	Fake *FakeEventingV1alpha1
	ns   string
}

var subscriptionsResource = schema.GroupVersionResource{Group: "eventing.kyma.cx", Version: "v1alpha1", Resource: "subscriptions"}

var subscriptionsKind = schema.GroupVersionKind{Group: "eventing.kyma.cx", Version: "v1alpha1", Kind: "Subscription"}

// Get takes name of the subscription, and returns the corresponding subscription object, and an error if there is any.
func (c *FakeSubscriptions) Get(name string, options v1.GetOptions) (result *v1alpha1.Subscription, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(subscriptionsResource, c.ns, name), &v1alpha1.Subscription{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Subscription), err
}

// List takes label and field selectors, and returns the list of Subscriptions that match those selectors.
func (c *FakeSubscriptions) List(opts v1.ListOptions) (result *v1alpha1.SubscriptionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(subscriptionsResource, subscriptionsKind, c.ns, opts), &v1alpha1.SubscriptionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SubscriptionList{}
	for _, item := range obj.(*v1alpha1.SubscriptionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested subscriptions.
func (c *FakeSubscriptions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(subscriptionsResource, c.ns, opts))

}

// Create takes the representation of a subscription and creates it.  Returns the server's representation of the subscription, and an error, if there is any.
func (c *FakeSubscriptions) Create(subscription *v1alpha1.Subscription) (result *v1alpha1.Subscription, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(subscriptionsResource, c.ns, subscription), &v1alpha1.Subscription{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Subscription), err
}

// Update takes the representation of a subscription and updates it. Returns the server's representation of the subscription, and an error, if there is any.
func (c *FakeSubscriptions) Update(subscription *v1alpha1.Subscription) (result *v1alpha1.Subscription, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(subscriptionsResource, c.ns, subscription), &v1alpha1.Subscription{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Subscription), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSubscriptions) UpdateStatus(subscription *v1alpha1.Subscription) (*v1alpha1.Subscription, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(subscriptionsResource, "status", c.ns, subscription), &v1alpha1.Subscription{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Subscription), err
}

// Delete takes name of the subscription and deletes it. Returns an error if one occurs.
func (c *FakeSubscriptions) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(subscriptionsResource, c.ns, name), &v1alpha1.Subscription{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSubscriptions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(subscriptionsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SubscriptionList{})
	return err
}

// Patch applies the patch and returns the patched subscription.
func (c *FakeSubscriptions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Subscription, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(subscriptionsResource, c.ns, name, data, subresources...), &v1alpha1.Subscription{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Subscription), err
}
