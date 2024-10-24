package comparator

import (
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/infrastructure-manager/hack/shoot-comparator/pkg/shoot"
	"github.com/pkg/errors"
)

type Result struct {
	Equal bool
	Diff  []Difference
}

type Difference struct {
	ShootName  string
	LeftShoot  v1beta1.Shoot
	RightShoot v1beta1.Shoot
	Message    string
}

func CompareShoots(leftShoot, rightShoot v1beta1.Shoot) (Result, error) {

	differences, err := compare(leftShoot, rightShoot)
	if err != nil {
		return Result{}, err
	}

	equal := len(differences) == 0

	return Result{
		Equal: equal,
		Diff:  differences,
	}, nil
}

func compare(leftShoot, rightShoot v1beta1.Shoot) ([]Difference, error) {
	var differences []Difference

	matcher := shoot.NewMatcher(leftShoot)
	equal, err := matcher.Match(rightShoot)
	if err != nil {
		return nil, err
	}

	if !equal {
		if leftShoot.Name != rightShoot.Name {
			return nil, errors.New("shoot names are different, stopping comparison")
		}
		diff := Difference{
			ShootName:  leftShoot.Name,
			LeftShoot:  leftShoot,
			RightShoot: rightShoot,
			Message:    matcher.FailureMessage(nil),
		}
		differences = append(differences, diff)
	}
	return differences, nil
}
