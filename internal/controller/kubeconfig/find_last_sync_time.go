package kubeconfig

import (
	"time"
)

func findLastSyncTime(annotations map[string]string) (time.Time, bool) {
	_, found := annotations[lastKubeconfigSyncAnnotation]
	if !found {
		return time.Time{}, false
	}

	lastSyncTimeString := annotations[lastKubeconfigSyncAnnotation]
	lastSyncTime, err := time.Parse(time.RFC3339, lastSyncTimeString)
	if err != nil {
		return time.Time{}, false
	}
	return lastSyncTime, true
}
