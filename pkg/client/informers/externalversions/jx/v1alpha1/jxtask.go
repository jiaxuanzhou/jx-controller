// Copyright 2018 2018 BY JIAXUAN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	jx_v1alpha1 "github.com/jiaxuanzhou/jx-controller/pkg/apis/jx/v1alpha1"
	versioned "github.com/jiaxuanzhou/jx-controller/pkg/client/clientset/versioned"
	internalinterfaces "github.com/jiaxuanzhou/jx-controller/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/jiaxuanzhou/jx-controller/pkg/client/listers/jx/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// JxTaskInformer provides access to a shared informer and lister for
// JxTasks.
type JxTaskInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.JxTaskLister
}

type jxTaskInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewJxTaskInformer constructs a new informer for JxTask type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewJxTaskInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredJxTaskInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredJxTaskInformer constructs a new informer for JxTask type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredJxTaskInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.JxV1alpha1().JxTasks(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.JxV1alpha1().JxTasks(namespace).Watch(options)
			},
		},
		&jx_v1alpha1.JxTask{},
		resyncPeriod,
		indexers,
	)
}

func (f *jxTaskInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredJxTaskInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *jxTaskInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&jx_v1alpha1.JxTask{}, f.defaultInformer)
}

func (f *jxTaskInformer) Lister() v1alpha1.JxTaskLister {
	return v1alpha1.NewJxTaskLister(f.Informer().GetIndexer())
}