package roundtrip

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"
	ujconversion "github.com/crossplane/upjet/v2/pkg/controller/conversion"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	fuzz "github.com/google/gofuzz"
	"github.com/hashicorp/terraform-provider-aws/xpprovider"
	clusterapis "github.com/upbound/provider-aws/v2/apis/cluster"
	"github.com/upbound/provider-aws/v2/config"
	"k8s.io/apimachinery/pkg/api/apitesting/fuzzer"
	"k8s.io/apimachinery/pkg/api/apitesting/roundtrip"
	genericfuzzer "k8s.io/apimachinery/pkg/apis/meta/fuzzer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/randfill"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/upbound/provider-aws/v2/apis/cluster/acmpca/v1beta1"
	"github.com/upbound/provider-aws/v2/apis/cluster/acmpca/v1beta2"
)

var (
	testScheme       *runtime.Scheme
	testCodecFactory serializer.CodecFactory
	setupOnce        sync.Once
	setupErr         error
)

// setupTestInfrastructure performs expensive one-time setup for all tests.
// It builds the scheme, retrieves the provider configuration, and registers conversions.
func setupTestInfrastructure() error {
	testScheme = runtime.NewScheme()
	if err := clusterapis.AddToScheme(testScheme); err != nil {
		return fmt.Errorf("failed to add cluster apis to scheme: %w", err)
	}

	testCodecFactory = serializer.NewCodecFactory(testScheme)

	fwProvider, sdkProvider, err := xpprovider.GetProvider(context.Background())
	if err != nil {
		return fmt.Errorf("cannot get the Terraform framework and SDK providers: %w", err)
	}

	p, err := config.GetProvider(context.TODO(), fwProvider, sdkProvider, false, false)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	if err := ujconversion.RegisterConversions(p, nil, testScheme); err != nil {
		return fmt.Errorf("failed to register conversions: %w", err)
	}

	return nil
}

// getTestInfrastructure returns the shared test infrastructure, initializing it once if needed.
// Returns the scheme and codec factory that are safe to read concurrently.
func getTestInfrastructure(t *testing.T) (*runtime.Scheme, serializer.CodecFactory) {
	t.Helper()
	setupOnce.Do(func() {
		setupErr = setupTestInfrastructure()
	})

	if setupErr != nil {
		t.Fatalf("test infrastructure setup failed: %v", setupErr)
	}

	return testScheme, testCodecFactory
}

// normalizeMeta strips fields that commonly differ and aren’t part of your API contract.
func normalizeMeta(obj metav1.Object) {
	// adjust based on your needs
	obj.SetResourceVersion("")
	obj.SetGeneration(0)
	obj.SetUID("")
	obj.SetManagedFields(nil)
	// timestamps are often nil in unit tests, but clear anyway if you set them somewhere
	obj.SetCreationTimestamp(metav1.Time{})
}

func cmpK8sObjects(t *testing.T, a, b runtime.Object) {
	t.Helper()

	// If your types implement metav1.Object, normalize metadata before comparing.
	if ao, ok := a.(interface{ GetObjectMeta() metav1.ObjectMeta }); ok {
		om := ao.GetObjectMeta()
		normalizeMeta(&om) // note: this copies; better to normalize via metav1.Object
	}
	// Prefer: cast to metav1.Object and mutate in-place
	if ao, ok := a.(metav1.Object); ok {
		normalizeMeta(ao)
	}
	if bo, ok := b.(metav1.Object); ok {
		normalizeMeta(bo)
	}

	// Add ignore rules for fields you intentionally default or reorder.
	opts := []cmp.Option{
		cmpopts.IgnoreMapEntries(func(k, v string) bool {
			return k == "internal.upjet.crossplane.io/field-conversions"
		}),
		cmpopts.EquateEmpty(),
		// If you have slices where order isn’t meaningful, add sorting opts here.
	}

	if diff := cmp.Diff(a, b, opts...); diff != "" {
		t.Fatalf("round-trip diff (-want +got):\n%s", diff)
	}
}

func TestMyKind_RoundTripLossless_Fuzzed(t *testing.T) {
	_, codecFactory := getTestInfrastructure(t)

	objFuzzer := fuzzer.FuzzerFor(
		fuzzer.MergeFuzzerFuncs(genericfuzzer.Funcs),
		rand.NewSource(rand.Int63()),
		codecFactory,
	).NumElements(0, 1).NilChance(0.5).MaxDepth(25)

	for i := 0; i < 10; i++ {
		orig := &v1beta1.CertificateAuthority{}
		objFuzzer.Fill(orig)
		orig.Namespace = ""

		hub := &v1beta2.CertificateAuthority{}
		// old -> hub
		if err := orig.ConvertTo(hub); err != nil {
			t.Fatalf("orig.ConvertTo(hub): %v", err)
		}

		// hub -> old
		back := &v1beta1.CertificateAuthority{}
		if err := back.ConvertFrom(hub); err != nil {
			t.Fatalf("back.ConvertFrom(hub): %v", err)
		}

		cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
	}
}

func TestSerialization_Roundtrip_Fuzzed(t *testing.T) {
	sche, codecFactory := getTestInfrastructure(t)

	objFuzzer := fuzzer.FuzzerFor(
		fuzzer.MergeFuzzerFuncs(genericfuzzer.Funcs, customFuzzFuncs),
		rand.NewSource(rand.Int63()),
		codecFactory,
	).NumElements(0, 5).NilChance(0.5).MaxDepth(25)

	roundtrip.RoundTripExternalTypesWithoutProtobuf(t, sche, codecFactory, objFuzzer, nil)
}

func TestAll_HubSpokeHub_Fuzzed(t *testing.T) {
	sche, codecFactory := getTestInfrastructure(t)

	groupKinds := groupsToKindFromScheme(sche)
	for group, kinds := range groupKinds {
		group := group // Capture for Go < 1.22
		t.Run(group, func(t *testing.T) {
			t.Parallel()

			// Create a fuzzer per parallel group to avoid race conditions on rand.Source
			groupFuzzer := fuzzer.FuzzerFor(
				fuzzer.MergeFuzzerFuncs(genericfuzzer.Funcs, customFuzzFuncs),
				rand.NewSource(rand.Int63()),
				codecFactory,
			).NumElements(0, 1).NilChance(0).MaxDepth(25)

			for _, gk := range kinds.UnsortedList() {
				gk := gk // Capture for Go < 1.22
				availableVersions := sche.VersionsForGroupKind(gk)
				if len(availableVersions) < 2 {
					t.Logf("Kind %q has one version, skipping conversion test", gk)
					continue
				}

				t.Run(gk.Kind, func(t *testing.T) {
					t.Parallel()
					testKind(t, sche, gk, availableVersions, groupFuzzer)
				})
			}
		})
	}
}

func customFuzzFuncs(_ serializer.CodecFactory) []interface{} {
	return []interface{}{
		// Force empty namespace for cluster-scoped resources
		func(meta *metav1.ObjectMeta, c randfill.Continue) {
			c.FillNoCustom(meta)
			meta.Namespace = "" // Always cluster-scoped
		},
	}
}

func testKind(t *testing.T, scheme *runtime.Scheme, gk schema.GroupKind, gvList []schema.GroupVersion, objFuzzer *randfill.Filler) {
	var hubVersion string
	spokes := make([]string, 0, len(gvList))
	for _, gv := range gvList {
		object, err := scheme.New(gk.WithVersion(gv.Version))
		if err != nil {
			t.Fatalf("cannot create object: %v", err)
		}
		if _, ok := object.(conversion.Hub); ok {
			hubVersion = gv.Version
		} else if _, ok := object.(conversion.Convertible); ok {
			spokes = append(spokes, gv.Version)
		}
	}

	if hubVersion == "" {
		t.Skipf("no hub version found for %s", gk)
	}
	if len(spokes) == 0 {
		t.Skipf("no spoke version found for %s", gk)
	}
	if len(spokes)+1 != len(gvList) {
		t.Skipf("missing spoke implementation for %s", gk)
	}
	for _, spokeVersion := range spokes {
		t.Run(fmt.Sprintf("%s_To_Hub_%s", spokeVersion, hubVersion), func(t *testing.T) {
			t.Parallel()
			spokeHubSpoke(t, scheme, gk.WithVersion(hubVersion), gk.WithVersion(spokeVersion), objFuzzer)
		})
		t.Run(fmt.Sprintf("%s_To_Spoke_%s", hubVersion, spokeVersion), func(t *testing.T) {
			t.Parallel()
			hubSpokeHub(t, scheme, gk.WithVersion(hubVersion), gk.WithVersion(spokeVersion), objFuzzer)
		})
	}
}

func spokeHubSpoke(t *testing.T, scheme *runtime.Scheme, hubGvk schema.GroupVersionKind, spokeGvk schema.GroupVersionKind, objFuzzer *randfill.Filler) {
	for i := 0; i < 10; i++ {
		origRuntime := resource.MustCreateObject(spokeGvk, scheme)
		orig, ok := origRuntime.(conversion.Convertible)
		if !ok {
			t.Fatalf("object is not Convertible: %v", origRuntime)
		}

		objFuzzer.Fill(origRuntime)

		hubRuntime := resource.MustCreateObject(hubGvk, scheme)
		hub, ok := hubRuntime.(conversion.Hub)
		if !ok {
			t.Fatalf("object is not Hub: %v", hubRuntime)
		}

		// old -> hub
		if err := orig.ConvertTo(hub); err != nil {
			t.Fatalf("orig.ConvertTo(hub): %v", err)
		}

		// hub -> old
		backRuntime := resource.MustCreateObject(spokeGvk, scheme)
		back, ok := backRuntime.(conversion.Convertible)
		if !ok {
			t.Fatalf("object is not Convertible: %v", backRuntime)
		}

		if err := back.ConvertFrom(hub); err != nil {
			t.Fatalf("back.ConvertFrom(hub): %v", err)
		}

		cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
	}
}

func hubSpokeHub(t *testing.T, scheme *runtime.Scheme, hubGvk schema.GroupVersionKind, spokeGvk schema.GroupVersionKind, objFuzzer *randfill.Filler) {
	for i := 0; i < 10; i++ {
		origRuntime := resource.MustCreateObject(hubGvk, scheme)
		orig, ok := origRuntime.(conversion.Hub)
		if !ok {
			t.Fatalf("object is not Convertible: %v", origRuntime)
		}

		objFuzzer.Fill(origRuntime)

		spokeRuntime := resource.MustCreateObject(spokeGvk, scheme)
		spoke, ok := spokeRuntime.(conversion.Convertible)
		if !ok {
			t.Fatalf("object is not Hub: %v", spokeRuntime)
		}

		// hub -> spoke
		if err := spoke.ConvertFrom(orig); err != nil {
			t.Fatalf("spoke.ConvertFrom(orig): %v", err)
		}

		backRuntime := resource.MustCreateObject(hubGvk, scheme)
		back, ok := backRuntime.(conversion.Hub)
		if !ok {
			t.Fatalf("object is not Convertible: %v", backRuntime)
		}

		// spoke -> hub
		if err := spoke.ConvertTo(back); err != nil {
			t.Fatalf("spoke.ConvertTo(back): %v", err)
		}

		cmpK8sObjects(t, orig.DeepCopyObject(), back.DeepCopyObject())
	}
}

var ignoreKinds = sets.New("ListOptions", "CreateOptions", "GetOptions", "UpdateOptions", "PatchOptions", "DeleteOptions", "WatchEvent", "PatchOptions")

func groupsToKindFromScheme(scheme *runtime.Scheme) map[string]sets.Set[schema.GroupKind] {
	ret := make(map[string]sets.Set[schema.GroupKind])

	for gvk := range scheme.AllKnownTypes() {
		if _, ok := ret[gvk.Group]; !ok {
			ret[gvk.Group] = sets.New[schema.GroupKind]()
		}
		if strings.HasSuffix(gvk.Kind, "List") || ignoreKinds.Has(gvk.Kind) {
			continue
		}
		ret[gvk.Group].Insert(gvk.GroupKind())
	}
	return ret
}

func EquateNilPtr() cmp.Option {
	return cmp.FilterPath(func(p cmp.Path) bool {
		return p.Last().Type().Kind() == reflect.Ptr
	}, cmp.Transformer("NilToZero", func(x interface{}) interface{} {
		v := reflect.ValueOf(x)
		if v.Kind() == reflect.Ptr && v.IsNil() {
			// Return a pointer to a zero value of the struct type
			return reflect.New(v.Type().Elem()).Interface()
		}
		return x
	}))
}

func NoNilElementsInPointerSlices(_ serializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj interface{}, c fuzz.Continue) {
			v := reflect.ValueOf(obj)
			if v.Kind() != reflect.Ptr {
				return
			}

			v = v.Elem()
			if v.Kind() != reflect.Slice {
				return
			}

			elemType := v.Type().Elem()
			if elemType.Kind() != reflect.Ptr {
				return
			}

			// Allow nil slice
			if c.RandBool() {
				v.Set(reflect.Zero(v.Type()))
				return
			}

			// Non-nil slice with no nil elements
			n := c.Intn(1) // tune as needed
			slice := reflect.MakeSlice(v.Type(), n, n)

			for i := 0; i < n; i++ {
				elem := reflect.New(elemType.Elem())
				c.Fuzz(elem.Interface())
				slice.Index(i).Set(elem)
			}

			v.Set(slice)
		},
	}
}
