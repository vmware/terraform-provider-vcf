# Private AI Foundation via Terraform - Samples

This repository contains sample automation for enabling a Tanzu Kubernetes Cluster with NVIDIA GPUs in VMware Cloud Foundation.

The configuration is divided into several steps which are intended to be executed in order.
Each step is designed to be atomic and can be executed independently from the rest provided that its prerequisites in terms of infrastructure are in place.

These examples also use the [vSphere Terraform Provider](https://github.com/hashicorp/terraform-provider-vsphere)

## Workflow Overview

### Starting State

This configuration is intended to be applied on an environment with a configured management domain.

### Desired State

Using these samples you can:

* Create a cluster image with vGPU drivers provided by NVIDIA
* Create a subscribed content library with custom container images
* Deploy a workload domain with vSAN storage and NSX network backing
* Create an NSX Edge Cluster
* Enable vSphere Supervisor on a cluster
* Configure a vSphere Namespace and a Virtual Machine Class with your vGPU drivers

## Contents

### The [steps](https://github.com/vmware/terraform-provider-vcf/-/tree/main/examples/workflows/private-ai-foundation/steps) directory contains the sample configuration divided into a number of examples/workflows/private-ai-foundation/steps
#### [Step 1](https://github.com/vmware/terraform-provider-vcf/-/blob/main/examples/workflows/private-ai-foundation/steps/01) - Create a vLCM cluster image with NVIDIA GPU drivers
#### [Step 2](https://github.com/vmware/terraform-provider-vcf/-/blob/main/examples/workflows/private-ai-foundation/steps/02) - Export the cluster image and create a workload domain with it
#### [Step 3](https://github.com/vmware/terraform-provider-vcf/-/blob/main/examples/workflows/private-ai-foundation/steps/03) - Create an NSX Edge Cluster
#### [Step 4](https://github.com/vmware/terraform-provider-vcf/-/blob/main/examples/workflows/private-ai-foundation/steps/04) - Create a subscribed Content Library
#### [Step 5](https://github.com/vmware/terraform-provider-vcf/-/blob/main/examples/workflows/private-ai-foundation/steps/05) - Enable Supervisor, create a vSphere Namespace and a VM class
