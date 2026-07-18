package controller

import (
	"fmt"
	"sort"

	nodev1alpha1 "brewlet-operator/api/nodeprofile/v1alpha1"
	"brewlet-operator/internal/brewlet"
)

// ValidateNodeProfile enforces the NodeProfile invariants the validating webhook
// rejects early (§8.3): a non-empty JDK list, only curated distributions, a
// valid containerdRestart value, and a positive feature version. The
// cross-profile "two profiles naming the same pool" check is
// ValidateNoPoolConflicts, which needs the sibling set.
func ValidateNodeProfile(profile *nodev1alpha1.NodeProfile) error {
	if len(profile.Spec.JDKs) == 0 {
		return fmt.Errorf("spec.jdks must list at least one JDK")
	}
	for i, j := range profile.Spec.JDKs {
		if !isCuratedDistribution(j.Distribution) {
			return fmt.Errorf("spec.jdks[%d].distribution %q is not curated; supported: %v",
				i, j.Distribution, brewlet.CuratedDistributions)
		}
		if j.Feature <= 0 {
			return fmt.Errorf("spec.jdks[%d].feature must be a positive version, got %d", i, j.Feature)
		}
	}
	switch profile.Spec.Rollout.ContainerdRestart {
	case "", nodev1alpha1.ContainerdRestartValidated,
		nodev1alpha1.ContainerdRestartSIGHUP, nodev1alpha1.ContainerdRestartNone:
	default:
		return fmt.Errorf("spec.rollout.containerdRestart %q is invalid; want one of validated|sighup|none",
			profile.Spec.Rollout.ContainerdRestart)
	}
	return nil
}

// ValidateNoPoolConflicts rejects a profile that names a pool already claimed by
// another profile (ambiguous ownership — the pool-model replacement for the old
// priority tie-break, §5.6/§8.3). others is the set of existing profiles
// (excluding the one under validation).
func ValidateNoPoolConflicts(profile *nodev1alpha1.NodeProfile, others []nodev1alpha1.NodeProfile) error {
	claimed := map[string]string{} // pool name -> owning profile
	for i := range others {
		o := &others[i]
		if o.Name == profile.Name {
			continue
		}
		for _, name := range o.Spec.NodePool.Names {
			claimed[name] = o.Name
		}
	}
	var dupes []string
	for _, name := range profile.Spec.NodePool.Names {
		if owner, ok := claimed[name]; ok {
			dupes = append(dupes, fmt.Sprintf("%q (already claimed by profile %q)", name, owner))
		}
	}
	if len(dupes) > 0 {
		sort.Strings(dupes)
		return fmt.Errorf("pool ownership conflict: %v", dupes)
	}
	return nil
}

func isCuratedDistribution(dist string) bool {
	for _, d := range brewlet.CuratedDistributions {
		if d == dist {
			return true
		}
	}
	return false
}
