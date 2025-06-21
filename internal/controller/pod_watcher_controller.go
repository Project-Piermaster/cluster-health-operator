package controller

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	healthv1alpha1 "github.com/Project-Piermaster/cluster-health-operator/api/v1alpha1"
)

type PodWatcher struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *PodWatcher) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var pod corev1.Pod
	if err := r.Get(ctx, req.NamespacedName, &pod); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if pod.Status.Phase != corev1.PodPending {
		return ctrl.Result{}, nil
	}

	// Skip if already has a node assigned
	if pod.Spec.NodeName != "" {
		return ctrl.Result{}, nil
	}

	// Add check to avoid flapping pods (wait X seconds in Pending)
	if time.Since(pod.CreationTimestamp.Time) < 30*time.Second {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
	}

	// TODO: Add scheduling evaluation logic
	// If no eligible nodes found → create ClusterIssue CR

	var nodeList corev1.NodeList
	if err := r.List(ctx, &nodeList); err != nil {
		return ctrl.Result{}, err
	}

	matchingNodes := GetSchedulableNodesForPod(&pod, nodeList.Items)

	if len(matchingNodes) == 0 {
		// pod is truly unschedulable, create ClusterIssue
		issue := &healthv1alpha1.ClusterIssue{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "pod-pending-",
			},
			Spec: healthv1alpha1.ClusterIssueSpec{
				PodRef: corev1.ObjectReference{
					Kind: "Pod", Name: pod.Name, Namespace: pod.Namespace, UID: pod.UID,
				},
				Reason:      "NoMatchingNode",
				Message:     "Pod does not match any available node due to taints, selectors, or readiness.",
				DiagnosedAt: metav1.Now(),
			},
		}
		// log the issue
		log.Printf("Creating ClusterIssue for pod %s in namespace %s", pod.Name, pod.Namespace)
		if err := r.Create(ctx, issue); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodWatcher) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Pod{}).
		Complete(r)
}
