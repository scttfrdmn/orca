// Package node provides Virtual Kubelet integration for ORCA.
//
// This package implements the bridge between ORCA's provider implementation
// and the Virtual Kubelet framework, allowing ORCA to register as a Kubernetes
// node and handle pod lifecycle events.
//
// The main components are:
// - Controller: Manages the Virtual Kubelet node lifecycle
// - VirtualKubeletAdapter: Adapts the ORCA provider to Virtual Kubelet interfaces
package node
