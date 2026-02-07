package roundtrip

import (
	"testing"

	v1beta1 "github.com/upbound/provider-aws/v2/apis/cluster/connect/v1beta1"
	v1beta2 "github.com/upbound/provider-aws/v2/apis/cluster/connect/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestRoutingProfile_RoundTrip_Over_Hub(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta1.RoutingProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta1.RoutingProfileSpec{
			ForProvider: v1beta1.RoutingProfileParameters{

				MediaConcurrencies: []v1beta1.MediaConcurrenciesParameters{
					{
						Channel: ptr.To("foo"),
					},
					{
						Channel: ptr.To("bar"),
					},
				},
			},
		},
		Status: v1beta1.RoutingProfileStatus{
			AtProvider: v1beta1.RoutingProfileObservation{
				QueueConfigsAssociated: []v1beta1.QueueConfigsAssociatedObservation{
					{
						QueueID: ptr.To("qc1"),
					},
					{
						QueueID: ptr.To("qc2"),
					},
				},
			},
		},
	}

	hub := &v1beta2.RoutingProfile{}
	if err := orig.ConvertTo(hub); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta1.RoutingProfile{}
	if err := back.ConvertFrom(hub); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}

func TestRoutingProfile_RoundTrip_Over_Spoke(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta2.RoutingProfile{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta2.RoutingProfileSpec{
			ForProvider: v1beta2.RoutingProfileParameters{
				MediaConcurrencies: []v1beta2.MediaConcurrenciesParameters{
					{
						Channel: ptr.To("foo"),
						CrossChannelBehavior: &v1beta2.CrossChannelBehaviorParameters{
							BehaviorType: ptr.To("foo-behavior"),
						},
					},
					{
						Channel: ptr.To("bar"),
						CrossChannelBehavior: &v1beta2.CrossChannelBehaviorParameters{
							BehaviorType: ptr.To("bar-behavior"),
						},
					},
				},
			},
		},
	}

	spoke := &v1beta1.RoutingProfile{}
	if err := spoke.ConvertFrom(orig); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta2.RoutingProfile{}
	if err := spoke.ConvertTo(back); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}
