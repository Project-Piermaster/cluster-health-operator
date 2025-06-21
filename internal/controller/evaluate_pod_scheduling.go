package controller

import (
	corev1 "k8s.io/api/core/v1"
)

// Returns a list of node names that *could* schedule this pod.
func GetSchedulableNodesForPod(pod *corev1.Pod, nodes []corev1.Node) []string {
	var matches []string

	for _, node := range nodes {
		if !isNodeReady(node) {
			continue
		}

		if !toleratesTaints(pod, node.Spec.Taints) {
			continue
		}

		if !matchesNodeSelector(pod, node) {
			continue
		}

		// Optionally check resource fit, affinity, etc.

		matches = append(matches, node.Name)
	}

	return matches
}

func isNodeReady(node corev1.Node) bool {
	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func toleratesTaints(pod *corev1.Pod, taints []corev1.Taint) bool {
	for _, taint := range taints {
		tolerated := false
		for _, toleration := range pod.Spec.Tolerations {
			if toleration.Key == taint.Key && toleration.Value == taint.Value && toleration.Effect == taint.Effect {
				tolerated = true
				break
			}
		}
		if !tolerated {
			return false
		}
	}
	return true
}

func matchesNodeSelector(pod *corev1.Pod, node corev1.Node) bool {
	selector := pod.Spec.NodeSelector
	if len(selector) == 0 {
		return true
	}

	for key, val := range selector {
		if nodeVal, ok := node.Labels[key]; !ok || nodeVal != val {
			return false
		}
	}

	return true
}
