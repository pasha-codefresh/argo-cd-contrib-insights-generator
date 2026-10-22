package pkg

import "testing"

func TestValidateRepo(t *testing.T) {
	for _, repo := range []string{"argoproj/argo-cd", "argoproj/argo-rollouts"} {
		if err := ValidateRepo(repo); err != nil {
			t.Errorf("ValidateRepo(%q) unexpected error: %v", repo, err)
		}
	}
	for _, repo := range []string{"", "argo-cd", "/argo-cd", "argoproj/", "a/b/c"} {
		if err := ValidateRepo(repo); err == nil {
			t.Errorf("ValidateRepo(%q) expected error", repo)
		}
	}
}
