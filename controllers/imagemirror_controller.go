/*
Copyright (c) 2020 Facebook, Inc. and its affiliates.

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
	"encoding/base64"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	slipwayk8sfacebookcomv1 "github.com/davidewatson/slipway/api/v1"
)

// ImageMirrorReconciler reconciles a ImageMirror object
type ImageMirrorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=slipway.k8s.facebook.com,resources=imagemirrors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slipway.k8s.facebook.com,resources=imagemirrors/status,verbs=get;update;patch

// Reconcile is called when a resource we are watching may have changed.
func (r *ImageMirrorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var imageMirror slipwayk8sfacebookcomv1.ImageMirror

	ctx := context.Background()
	log := r.Log.WithValues("imagemirror", req.NamespacedName)

	// Get current version of the spec.
	if err := r.Get(ctx, req.NamespacedName, &imageMirror); err != nil {
		log.Error(err, "unable to fetch ImageMirror")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get credentials needed to mirror.
	sourceToken, err := r.GetToken(ctx, imageMirror.ObjectMeta.Namespace, imageMirror.Spec.SourceSecretName)
	if err != nil {
		log.Error(err, "unable to GetToken for source")
		return ctrl.Result{}, err
	}
	destToken, err := r.GetToken(ctx, imageMirror.ObjectMeta.Namespace, imageMirror.Spec.DestSecretName)
	if err != nil {
		log.Error(err, "unable to GetToken for dest")
		return ctrl.Result{}, err
	}

	// Mirror tags based on the users intent.
	mirroredTags, err := MirrorImage(ctx, imageMirror, log, sourceToken, destToken)
	if err != nil {
		return ctrl.Result{RequeueAfter: time.Minute}, err
	}

	// Update status with the current state.
	imageMirror.Status.MirroredTags = mirroredTags
	if err := r.Status().Update(ctx, &imageMirror); err != nil {
		log.Error(err, "unable to update ImageMirror status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// GetToken returns the token contained within the secret named name in namespace.
func (r *ImageMirrorReconciler) GetToken(ctx context.Context, namespace, name string) (string, error) {
	// Get the resource using a typed object.
	secret := &corev1.Secret{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, secret); err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(string(secret.Data["token"]))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// SetupWithManager registers controller with manager and configures shared informer.
func (r *ImageMirrorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slipwayk8sfacebookcomv1.ImageMirror{}).
		Complete(r)
}
