/* eslint-disable global-require */

import {sortBy} from '@site/src/utils/jsUtils';

/*
 * ADD YOUR PACKAGE TO THE GLASSKUBE PACKAGE OVERVIEW
 *
 * Please note that the packages displayed on the websites also include coming soon / planed packages.
 * Fore more information join the discussion on GitHub:
 *
 * https://github.com/glasskube/glasskube/discussions/90
 *
 */

export type TagType =
  | 'ai'
  | 'backup'
  | 'configuration'
  | 'planned'
  | 'database'
  | 'delivery'
  | 'logging'
  | 'messaging'
  | 'monitoring'
  | 'networking'
  | 'security'
  | 'visualization'
  | 'webassembly';

// Add sites to this list
// prettier-ignore
const Users: Package[] = [
  {
    name: 'Akri',
    shortDescription: 'Akri lets you easily expose heterogeneous leaf devices (such as IP cameras and USB devices) as resources in a Kubernetes cluster',
    iconUrl: 'https://avatars.githubusercontent.com/u/91915613',
    websiteUrl: 'https://docs.akri.sh/',
    sourceUrl: 'https://github.com/project-akri/akri',
    tags: ['configuration'],
  },
  {
    name: 'Argo CD',
    shortDescription: 'Declarative Continuous Deployment for Kubernetes',
    iconUrl: 'https://avatars.githubusercontent.com/u/30269780',
    websiteUrl: 'https://argo-cd.readthedocs.io',
    sourceUrl: 'https://github.com/argoproj/argo-cd',
    tags: ['configuration', 'delivery', 'visualization'],
  },
  {
    name: 'Argo Events',
    shortDescription: 'Argo Events is an event-driven workflow automation framework for Kubernetes.',
    iconUrl: 'https://avatars.githubusercontent.com/u/30269780',
    websiteUrl: 'https://argoproj.github.io/events/',
    sourceUrl: 'https://github.com/argoproj/argo-events',
    tags: ['planned', 'configuration', 'delivery', 'visualization'],
  },
  {
    name: 'Argo Rollouts',
    shortDescription: 'Argo Rollouts is a Kubernetes controller and set of CRDs which provides advanced deployment capabilities.',
    iconUrl: 'https://avatars.githubusercontent.com/u/30269780',
    websiteUrl: 'https://argoproj.github.io/rollouts/',
    sourceUrl: 'https://github.com/argoproj/argo-rollouts',
    tags: ['planned', 'configuration', 'delivery', 'visualization'],
  },
  {
    name: 'Argo Workflows',
    shortDescription: 'Kubernetes-native workflow engine supporting DAG and step-based workflows.',
    iconUrl: 'https://avatars.githubusercontent.com/u/30269780',
    websiteUrl: 'https://argoproj.github.io/workflows/',
    sourceUrl: 'https://github.com/argoproj/argo-workflows',
    tags: ['planned', 'configuration', 'delivery', 'visualization'],
  },
  {
    name: 'Caddy Ingress Controller',
    shortDescription: 'This is the Kubernetes Ingress Controller for Caddy',
    iconUrl: 'https://avatars.githubusercontent.com/u/12955528',
    websiteUrl: 'https://caddyserver.com',
    sourceUrl: 'https://github.com/caddyserver/ingress',
    tags: ['networking'],
  },
  {
    name: 'Cert manager',
    shortDescription: 'Automatically provision and manage TLS certificates in Kubernetes',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/548c320d-9b39-4578-9ddc-76f6e385ffbf',
    websiteUrl: 'https://cert-manager.io/',
    sourceUrl: 'https://github.com/cert-manager/cert-manager',
    tags: ['security'],
  },
  {
    name: 'CloudNativePG',
    shortDescription: 'CloudNativePG covers the full lifecycle of a PostgreSQL database cluster',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/fed38d57-1a62-44a3-bcb2-f07c9d6ab803',
    websiteUrl: 'https://cloudnative-pg.io/',
    sourceUrl: 'https://github.com/cloudnative-pg/cloudnative-pg',
    tags: ['database'],
  },
  {
    name: 'Cyclops',
    shortDescription: 'Developer friendly Kubernetes 👁️',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/8ff6069e-f481-4b8d-9a43-4cae6730d1ac',
    websiteUrl: 'https://cyclops-ui.com/',
    sourceUrl: 'https://github.com/cyclops-ui/cyclops',
    tags: ['visualization', 'configuration'],
  },
  {
    name: 'Gateway api',
    shortDescription: 'The Gateway API provides a standardized way to configure and manage network traffic routing in Kubernetes clusters.',
    iconUrl: 'https://avatars.githubusercontent.com/u/36015203',
    websiteUrl: 'https://gateway-api.sigs.k8s.io/',
    sourceUrl: 'https://github.com/kubernetes-sigs/gateway-api',
    tags: ['networking'],
  },
  {
    name: 'GPU Operator',
    shortDescription: 'The NVIDIA GPU Operator automates the management of all NVIDIA software components needed to provision GPU.',
    iconUrl: 'https://avatars.githubusercontent.com/u/1728152',
    websiteUrl: 'https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/index.html',
    sourceUrl: 'https://github.com/NVIDIA/gpu-operator',
    tags: ['planned', 'configuration'],
  },
  {
    name: 'Grafana',
    shortDescription: 'Use Grafana to visualize Metrics you collected in your Kubernetes cluster',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/3a7d6bbf-eb54-4353-88c1-624acb73b8aa',
    websiteUrl: 'https://prometheus-operator.dev/',
    sourceUrl: 'https://github.com/prometheus-operator/kube-prometheus',
    tags: ['monitoring'],
  },
  {
    name: 'Hatchet',
    shortDescription: 'A distributed, fault-tolerant task queue',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/7a01a56c-dc74-4afc-90af-17f79c3077d1',
    websiteUrl: 'https://hatchet.run/',
    sourceUrl: 'https://github.com/hatchet-dev/hatchet',
    tags: ['planned', 'visualization', 'messaging'],
  },
  {
    name: 'Headlamp',
    shortDescription: 'A Kubernetes web UI that is fully-featured, user-friendly and extensible',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/99016aa3-7033-4568-ad3c-b5e6e74105ad',
    websiteUrl: 'https://headlamp.dev/',
    sourceUrl: 'https://github.com/headlamp-k8s/headlamp',
    tags: ['visualization', 'configuration'],
  },
  {
    name: 'ingress-nginx',
    shortDescription: 'Ingress-NGINX Controller for Kubernetes',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/6701cd57-b690-4641-b967-ef2faee646e5',
    websiteUrl: 'https://kubernetes.github.io/ingress-nginx/',
    sourceUrl: 'https://github.com/kubernetes/ingress-nginx',
    tags: ['networking'],
  },
  {
    name: 'Kubernetes Dashboard',
    shortDescription: 'General-purpose web UI for Kubernetes clusters',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/6701cd57-b690-4641-b967-ef2faee646e5',
    websiteUrl: 'https://github.com/kubernetes/dashboard',
    sourceUrl: 'https://github.com/kubernetes/dashboard',
    tags: ['visualization', 'configuration'],
  },
  {
    name: 'Keptn',
    shortDescription: 'Toolkit for cloud-native application lifecycle management',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/40fc3a08-d9c6-4c00-ac5f-3db22411443b',
    websiteUrl: 'https://keptn.sh/',
    sourceUrl: 'https://github.com/keptn/lifecycle-toolkit',
    tags: ['delivery', 'logging'],
  },
  {
    name: 'Kubetail',
    shortDescription: 'Kubetail is a web-based, real-time log viewer for Kubernetes clusters',
    iconUrl: 'https://avatars.githubusercontent.com/u/141319696',
    websiteUrl: 'https://github.com/kubetail-org/kubetail',
    sourceUrl: 'https://github.com/kubetail-org/kubetail',
    tags: ['visualization', 'logging'],
  },
  {
    name: 'K8sGPT',
    shortDescription: 'Automatic SRE Superpowers within your Kubernetes cluster',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/3aee51a5-acb5-4575-b4fe-0e3df6591c83',
    websiteUrl: 'https://k8sgpt.ai/',
    sourceUrl: 'https://github.com/k8sgpt-ai/k8sgpt-operator',
    tags: ['ai'],
  },
  {
    name: 'Kube Prometheus Stack',
    shortDescription: 'Use Prometheus to collect Metrics from applications running on Kubernetes',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/afd860f4-d035-48e3-8f2d-36c632b9ff78',
    websiteUrl: 'https://prometheus-operator.dev/',
    sourceUrl: 'https://github.com/prometheus-operator/kube-prometheus',
    tags: ['monitoring'],
  },
  {
    name: 'Kubeflow',
    shortDescription: 'RabbitMQ Cluster Kubernetes Operator',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/821e6f42-472d-453e-8fb2-c0295c9a17f6',
    websiteUrl: 'https://www.kubeflow.org/',
    sourceUrl: 'https://github.com/kubeflow/kubeflow',
    tags: ['planned', 'ai', 'visualization', 'configuration', 'delivery'],
  },
  {
    name: 'Kwasm',
    shortDescription: 'KWasm allows fine-grained control over node provisioning',
    iconUrl: 'https://avatars.githubusercontent.com/u/113929551',
    websiteUrl: 'https://kwasm.sh/',
    sourceUrl: 'https://github.com/KWasm/kwasm-operator/',
    tags: ['webassembly', 'delivery', 'planned'],
  },
  {
    name: 'Litmus',
    shortDescription: 'Litmus helps SREs and developers practice chaos engineering in a Cloud-native way',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/7e791499-5677-476b-8cd4-5172ac235e5a',
    websiteUrl: 'https://litmuschaos.io/',
    sourceUrl: 'https://github.com/litmuschaos/litmus',
    tags: ['planned', 'security'],
  },
  {
    name: 'MariaDB Operator',
    shortDescription: 'Run and operate MariaDB in a cloud native way',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/12c90857-6a3c-416c-bfac-69a26e212240',
    websiteUrl: 'https://github.com/mariadb-operator/mariadb-operator',
    sourceUrl: 'https://github.com/mariadb-operator/mariadb-operator',
    tags: ['planned', 'database'],
  },
  {
    name: 'Node Feature Discovery',
    shortDescription: 'Node Feature Discovery (NFD) is a Kubernetes add-on for detecting hardware features and system configuration.',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/6701cd57-b690-4641-b967-ef2faee646e5',
    websiteUrl: 'https://github.com/NVIDIA/gpu-operator/tree/main/deployments/gpu-operator/charts/node-feature-discovery',
    sourceUrl: 'https://github.com/NVIDIA/gpu-operator/tree/main/deployments/gpu-operator/charts/node-feature-discovery',
    tags: ['configuration'],
  },
  {
    name: 'Quickwit',
    shortDescription: 'Cloud-native search engine for observability',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/16959694/c36ac826-6f80-4e6b-9a08-9719759e94b6',
    websiteUrl: 'https://quickwit.io/',
    sourceUrl: 'https://github.com/quickwit-oss/quickwit',
    tags: ['logging'],
  },
  {
    name: 'RabbitMQ Operator',
    shortDescription: 'RabbitMQ Cluster Kubernetes Operator',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/adfe6b1b-625c-4344-aecb-416cd7fccea7',
    websiteUrl: 'https://www.rabbitmq.com/kubernetes/operator/operator-overview',
    sourceUrl: 'https://github.com/rabbitmq/cluster-operator',
    tags: ['messaging'],
  },
  {
    name: 'SpinKube',
    shortDescription: 'The Spin Operator enables deploying Spin applications to Kubernetes.',
    iconUrl: 'https://avatars.githubusercontent.com/u/157650797',
    websiteUrl: 'https://www.spinkube.dev/',
    sourceUrl: 'https://github.com/spinkube/spin-operator',
    tags: ['webassembly', 'delivery', 'planned'],
  },
  {
    name: 'Velero',
    shortDescription: 'Backup and migrate Kubernetes applications and their persistent volumes',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/943fef9d-629c-4a93-87a7-45193c822fb2',
    websiteUrl: 'https://velero.io/',
    sourceUrl: 'https://github.com/vmware-tanzu/velero',
    tags: ['planned', 'backup'],
  }
  /*
  * Pro Tip: add your site in alphabetical order.
  * Appending your site here (at the end) is more likely to produce Git conflicts.
   */
];

export type Package = {
  name: string;
  shortDescription: string;
  iconUrl: string;
  websiteUrl: string;
  sourceUrl: string;
  tags: TagType[];
};

export type Tag = {
  label: string;
  color: string;
};

export const Tags: {[type in TagType]: Tag} = {
  ai: {
    label: 'ai',
    color: '#39ca30',
  },

  backup: {
    label: 'backup',
    color: '#dfd545',
  },

  configuration: {
    label: 'configuration',
    color: '#a44fb7',
  },

  planned: {
    label: 'coming soon',
    color: '#127f82',
  },

  database: {
    label: 'database',
    color: '#fe6829',
  },

  delivery: {
    label: 'delivery',
    color: '#8c2f00',
  },

  logging: {
    label: 'logging',
    color: '#1555da',
  },

  messaging: {
    label: 'messaging',
    color: '#cf6814',
  },

  monitoring: {
    label: 'monitoring',
    color: '#14cfc3',
  },

  networking: {
    label: 'networking',
    color: '#ffcfc3',
  },

  security: {
    label: 'security',
    color: '#a32cab',
  },

  visualization: {
    label: 'visualization',
    color: '#ab9a2c',
  },

  webassembly: {
    label: 'webassembly',
    color: '#8b6a2f',
  },
};

export const TagList = Object.keys(Tags) as TagType[];

function sortUsers() {
  let result = Users;
  // Sort by site name
  result = sortBy(result, user => user.name.toLowerCase());
  // Sort by favorite tag, favorites first
  result = sortBy(result, user => user.tags.includes('planned'));
  return result;
}

export const sortedUsers = sortUsers();
