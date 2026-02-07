package roundtrip

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	v1beta1 "github.com/upbound/provider-aws/v2/apis/cluster/redshift/v1beta1"
	v1beta2 "github.com/upbound/provider-aws/v2/apis/cluster/redshift/v1beta2"
)

func TestMyKind_RoundTrip_Over_Hub(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta2.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta2.ClusterSpec{
			ForProvider: v1beta2.ClusterParameters{
				Endpoint:         ptr.To("fooep"),
				ApplyImmediately: ptr.To(false),
				ClusterPublicKey: ptr.To("foopub"),
			},
		},
	}

	hub := &v1beta1.Cluster{}
	if err := orig.ConvertTo(hub); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta2.Cluster{}
	if err := back.ConvertFrom(hub); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}

func TestMyKind_RoundTrip_Over_Spoke(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta1.ClusterSpec{
			ForProvider: v1beta1.ClusterParameters{
				Encrypted: pointer.String("true"),
			},
		},
	}

	spoke := &v1beta2.Cluster{}
	if err := spoke.ConvertFrom(orig); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta1.Cluster{}
	if err := spoke.ConvertTo(back); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}

//func TestMyKind_RoundTripLossless(t *testing.T) {
//	orig := &v1beta1.Cluster{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      "example",
//			Namespace: "default",
//			Labels: map[string]string{
//				"app": "demo",
//			},
//			Annotations: map[string]string{
//				// include any “preserve unknown fields” annotations you use
//			},
//		},
//		Spec: v1beta1.ClusterSpec{
//			// Fill ALL interesting fields, including edge cases:
//			// - optional pointers
//			// - oneOf-ish structs
//			// - maps with multiple keys
//			// - nested lists
//		},
//	}
//
//	spoke := &v1beta2.Cluster{}
//
//	final := roundTrip(t, orig, spoke).(*v1beta1.Cluster)
//	cmpK8sObjects(t, orig.DeepCopyObject(), final.DeepCopyObject())
//}
/*
orig := &v1beta1.CertificateAuthority{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta1.CertificateAuthoritySpec{
			ForProvider: v1beta1.CertificateAuthorityParameters{
				CertificateAuthorityConfiguration: []v1beta1.CertificateAuthorityConfigurationParameters{
					{
						SigningAlgorithm: ptr.String("ECDSA"),
						Subject:          nil,
					},
				},
			},
		},
	}
*/
// roundTrip converts hub -> spoke -> hub and returns the final hub.
func roundTrip(t *testing.T, hub conversion.Hub, spoke conversion.Convertible) conversion.Hub {
	t.Helper()

	// hub -> spoke
	{
		dst := spoke.DeepCopyObject().(conversion.Convertible)
		if err := dst.ConvertFrom(hub); err != nil {
			t.Fatalf("ConvertFrom(hub) failed: %v", err)
		}
		spoke = dst
	}

	// spoke -> hub
	{
		dst := hub.DeepCopyObject().(conversion.Hub)
		if err := spoke.ConvertTo(dst); err != nil {
			t.Fatalf("ConvertTo(hub) failed: %v", err)
		}
		hub = dst
	}

	return hub
}
