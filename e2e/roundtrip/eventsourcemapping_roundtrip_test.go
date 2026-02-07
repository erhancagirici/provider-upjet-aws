package roundtrip

import (
	"testing"

	v1beta1 "github.com/upbound/provider-aws/v2/apis/cluster/lambda/v1beta1"
	v1beta2 "github.com/upbound/provider-aws/v2/apis/cluster/lambda/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func Test_Wrapped(t *testing.T) {
	_, _ = getTestInfrastructure(t)
	t.Run("Test_RT_Over_Hub", func(t *testing.T) {
		TestEventSourceMapping_RoundTrip_Over_Hub(t)
	})
	t.Run("Test_RT_Over_Spoke", func(t *testing.T) {
		TestEventSourceMapping_RoundTrip_Over_Spoke(t)
	})
}

func TestEventSourceMapping_RoundTrip_Over_Hub(t *testing.T) {
	_, _ = getTestInfrastructure(t)
	//
	//orig := &v1beta1.EventSourceMapping{
	//	ObjectMeta: metav1.ObjectMeta{
	//		Name: "example",
	//	},
	//	Spec: v1beta1.EndpointSpec{
	//		ForProvider: v1beta1.EndpointParameters{
	//			DatabaseName: ptr.To("example"),
	//		},
	//	},
	//}
	//
	//hub := &v1beta2.EventSourceMapping{}
	//if err := orig.ConvertTo(hub); err != nil {
	//	t.Fatalf("orig.ConvertTo(hub): %v", err)
	//}
	//
	//back := &v1beta1.EventSourceMapping{}
	//if err := back.ConvertFrom(hub); err != nil {
	//	t.Fatalf("back.ConvertFrom(hub): %v", err)
	//}
	//
	//cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}

func TestEventSourceMapping_RoundTrip_Over_Spoke(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta2.EventSourceMapping{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta2.EventSourceMappingSpec{
			ForProvider: v1beta2.EventSourceMappingParameters{
				AmazonManagedKafkaEventSourceConfig: &v1beta2.AmazonManagedKafkaEventSourceConfigParameters{
					ConsumerGroupID: ptr.To("ɥ(蚕ȔHQa¿湦蕼囿樾1"),
					SchemaRegistryConfig: &v1beta2.SchemaRegistryConfigParameters{
						AccessConfig: []v1beta2.AccessConfigParameters{
							{
								Type: ptr.To("L"),
								URI:  ptr.To("a"),
							},
						},
						EventRecordFormat:      ptr.To("_Ƀ·.鈪R"),
						SchemaRegistryURI:      ptr.To(`胚彩妥ʭ埛ZÅ硅ŰH\蓘@ȚF5`),
						SchemaValidationConfig: []v1beta2.SchemaValidationConfigParameters{},
					},
				},
			},
		},
	}

	spoke := &v1beta1.EventSourceMapping{}
	if err := spoke.ConvertFrom(orig); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta2.EventSourceMapping{}
	if err := spoke.ConvertTo(back); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}
