package controller

import (
	restclientset "k8s.io/client-go/rest"

	"github.com/jiaxuanzhou/jx-controller/pkg/client/informers/externalversions/jx/v1alpha1"
)

func NewJxTaskInformer(restConfig *restclientset.Config) v1alpha1.JxTaskInformer {

	return nil
}