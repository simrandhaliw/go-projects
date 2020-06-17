// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	templatev1 "github.com/openshift/api/template/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeBrokerTemplateInstances implements BrokerTemplateInstanceInterface
type FakeBrokerTemplateInstances struct {
	Fake *FakeTemplateV1
}

var brokertemplateinstancesResource = schema.GroupVersionResource{Group: "template.openshift.io", Version: "v1", Resource: "brokertemplateinstances"}

var brokertemplateinstancesKind = schema.GroupVersionKind{Group: "template.openshift.io", Version: "v1", Kind: "BrokerTemplateInstance"}

// Get takes name of the brokerTemplateInstance, and returns the corresponding brokerTemplateInstance object, and an error if there is any.
func (c *FakeBrokerTemplateInstances) Get(ctx context.Context, name string, options v1.GetOptions) (result *templatev1.BrokerTemplateInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(brokertemplateinstancesResource, name), &templatev1.BrokerTemplateInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*templatev1.BrokerTemplateInstance), err
}

// List takes label and field selectors, and returns the list of BrokerTemplateInstances that match those selectors.
func (c *FakeBrokerTemplateInstances) List(ctx context.Context, opts v1.ListOptions) (result *templatev1.BrokerTemplateInstanceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(brokertemplateinstancesResource, brokertemplateinstancesKind, opts), &templatev1.BrokerTemplateInstanceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &templatev1.BrokerTemplateInstanceList{ListMeta: obj.(*templatev1.BrokerTemplateInstanceList).ListMeta}
	for _, item := range obj.(*templatev1.BrokerTemplateInstanceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested brokerTemplateInstances.
func (c *FakeBrokerTemplateInstances) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(brokertemplateinstancesResource, opts))
}

// Create takes the representation of a brokerTemplateInstance and creates it.  Returns the server's representation of the brokerTemplateInstance, and an error, if there is any.
func (c *FakeBrokerTemplateInstances) Create(ctx context.Context, brokerTemplateInstance *templatev1.BrokerTemplateInstance, opts v1.CreateOptions) (result *templatev1.BrokerTemplateInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(brokertemplateinstancesResource, brokerTemplateInstance), &templatev1.BrokerTemplateInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*templatev1.BrokerTemplateInstance), err
}

// Update takes the representation of a brokerTemplateInstance and updates it. Returns the server's representation of the brokerTemplateInstance, and an error, if there is any.
func (c *FakeBrokerTemplateInstances) Update(ctx context.Context, brokerTemplateInstance *templatev1.BrokerTemplateInstance, opts v1.UpdateOptions) (result *templatev1.BrokerTemplateInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(brokertemplateinstancesResource, brokerTemplateInstance), &templatev1.BrokerTemplateInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*templatev1.BrokerTemplateInstance), err
}

// Delete takes name of the brokerTemplateInstance and deletes it. Returns an error if one occurs.
func (c *FakeBrokerTemplateInstances) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(brokertemplateinstancesResource, name), &templatev1.BrokerTemplateInstance{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeBrokerTemplateInstances) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(brokertemplateinstancesResource, listOpts)

	_, err := c.Fake.Invokes(action, &templatev1.BrokerTemplateInstanceList{})
	return err
}

// Patch applies the patch and returns the patched brokerTemplateInstance.
func (c *FakeBrokerTemplateInstances) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *templatev1.BrokerTemplateInstance, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(brokertemplateinstancesResource, name, pt, data, subresources...), &templatev1.BrokerTemplateInstance{})
	if obj == nil {
		return nil, err
	}
	return obj.(*templatev1.BrokerTemplateInstance), err
}