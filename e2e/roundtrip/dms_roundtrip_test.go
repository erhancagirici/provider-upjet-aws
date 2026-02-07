package roundtrip

import (
	"testing"

	v1beta1 "github.com/upbound/provider-aws/v2/apis/cluster/dms/v1beta1"
	v1beta2 "github.com/upbound/provider-aws/v2/apis/cluster/dms/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func TestDMS_RoundTrip_Over_Hub(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta1.Endpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta1.EndpointSpec{
			ForProvider: v1beta1.EndpointParameters{
				DatabaseName: ptr.To("example"),
			},
		},
	}

	hub := &v1beta2.Endpoint{}
	if err := orig.ConvertTo(hub); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta1.Endpoint{}
	if err := back.ConvertFrom(hub); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}

func TestDMS_RoundTrip_Over_Spoke(t *testing.T) {
	_, _ = getTestInfrastructure(t)

	orig := &v1beta2.Endpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: v1beta2.EndpointSpec{
			ForProvider: v1beta2.EndpointParameters{
				DatabaseName: ptr.To("example"),
				OracleSettings: &v1beta2.OracleSettingsParameters{
					AuthenticationMethod: ptr.To("foo"),
				},
				//MySQLSettings: &v1beta2.MySQLSettingsParameters{
				//	AfterConnectScript:            ptr.To("ʒ钒ŀ"),
				//	AuthenticationMethod:          ptr.To("ʣĜʉ冞 "),
				//	CleanSourceMetadataOnMismatch: ptr.To(false),
				//	EventsPollInterval:            ptr.To(0.051643313320298805),
				//	ExecuteTimeout:                ptr.To(0.9936911406785721),
				//	MaxFileSize:                   ptr.To(0.35230656228922325),
				//	ParallelLoadThreads:           ptr.To(0.8642498654858927),
				//	ServerTimezone:                ptr.To("Ò@ɑúōĚ茊癩庌:佷Ʌ霸崭Ȑ还"),
				//	ServiceAccessRoleArn:          ptr.To("ś~鳟/曀w剔種ʈ垏訾Ǭƃǜȵʪ"),
				//	TargetDBType:                  ptr.To("!ĻO"),
				//},
			},
		},
	}

	spoke := &v1beta1.Endpoint{}
	if err := spoke.ConvertFrom(orig); err != nil {
		t.Fatalf("orig.ConvertTo(hub): %v", err)
	}

	back := &v1beta2.Endpoint{}
	if err := spoke.ConvertTo(back); err != nil {
		t.Fatalf("back.ConvertFrom(hub): %v", err)
	}

	cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
}
