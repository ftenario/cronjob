/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

   "github.com/go-logr/logr"
    "github.com/robfig/cron"
    kbatch "k8s.io/api/batch/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "tutorial.kuberbuilder.io/cronjob/api/v1"
)

// CronJobReconciler reconciles a CronJob object
type CronJobReconciler struct {
	client.Client
	Log logr.Logger
	Scheme *runtime.Scheme
	Clock
}

type realClock struct{}

func (_ realClock) Now() time.Time { return time.Now() }

// clock knows how to get the current time.
// It can be used to fake out timing for testing.
type Clock interface {
    Now() time.Time
}

var (
    scheduledTimeAnnotation = "batch.tutorial.kubebuilder.io/scheduled-at"
)

//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=cronjobs/finalizers,verbs=update
//+kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CronJob object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *CronJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WihValues("cronjob", req.NamespaceName)

	// your logic here

	//Load the CronJob by Name
	var cronJob batchv1.Cronjob
	if err := r.Get(ctx, req.NamespaceName, &cronJob); err != nil {
		log.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they cant be fixed by an immediate
		// requeue (well need to wait for a new notification), and we can get them on deleted requests
		return crtl.Result{}, client.IgnoreNotFound(err) 
	}

	//List all active jobs, and update the status
	var childJobs kbatch.JobList
	if err := r.List(ctx, &childJobs, client.InNameSpace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error("err", unable to list child jobs)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CronJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.CronJob{}).
		Complete(r)
}
