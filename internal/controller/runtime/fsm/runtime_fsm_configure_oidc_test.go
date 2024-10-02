package fsm

import (
	"context"
	"fmt"
	gardener "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"testing"
)

func TestOidcState(t *testing.T) {
	t.Run("Should switch state to ApplyClusterRoleBindings when extension is disabled", func(t *testing.T) {
		// given
		ctx := context.Background()
		fsm := &fsm{}

		runtimeStub := runtimeForTest()
		shootStub := shootForTest()
		oidcService := gardener.Extension{
			Type:     "OidcExtensionType",
			Disabled: ptr.To(true),
		}
		shootStub.Spec.Extensions = append(shootStub.Spec.Extensions, oidcService)

		systemState := &systemState{
			instance: runtimeStub,
			shoot:    shootStub,
		}

		expectedRuntimeConditions := []metav1.Condition{
			{
				Type:    string(imv1.ConditionTypeOidcConfigured),
				Reason:  string(imv1.ConditionReasonOidcConfigured),
				Status:  "True",
				Message: "OIDC extension disabled",
			},
		}

		//when
		stateFn, _, _ := sFnConfigureOidc(ctx, fsm, systemState)

		//then
		require.Contains(t, stateFn.name(), "sFnApplyClusterRoleBindings")
		assertEqualConditions(t, expectedRuntimeConditions, systemState.instance.Status.Conditions)
	})

}

// sets the time to its zero value for comparison purposes
func assertEqualConditions(t *testing.T, expectedConditions []metav1.Condition, actualConditions []metav1.Condition, msgAndArgs ...interface{}) bool {
	for i, _ := range actualConditions {
		time := metav1.Time{}
		actualConditions[i].LastTransitionTime = time
		fmt.Print("a")
	}

	return assert.Equal(t, expectedConditions, actualConditions)
}
