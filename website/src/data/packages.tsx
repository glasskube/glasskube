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
  | 'monitoring'
  | 'networking'
  | 'security'
  | 'visualization'
  ;

// Add sites to this list
// prettier-ignore
const Users: Package[] = [
  {
    name: 'Cert manager',
    shortDescription: 'Automatically provision and manage TLS certificates in Kubernetes',
    iconUrl: 'https://avatars.githubusercontent.com/u/39950598?s=150',
    websiteUrl: 'https://cert-manager.io/',
    sourceUrl: 'https://github.com/cert-manager/cert-manager',
    tags: ['security'],
  },
  {
    name: 'ingress-nginx',
    shortDescription: 'Ingress-NGINX Controller for Kubernetes',
    iconUrl: 'https://avatars.githubusercontent.com/u/1412239?s=150',
    websiteUrl: 'https://kubernetes.github.io/ingress-nginx/',
    sourceUrl: 'https://github.com/kubernetes/ingress-nginx',
    tags: ['networking'],
  },
  {
    name: 'Kubernetes Dashboard',
    shortDescription: 'General-purpose web UI for Kubernetes clusters',
    iconUrl: 'https://avatars.githubusercontent.com/u/13629408?s=150',
    websiteUrl: 'https://github.com/kubernetes/dashboard',
    sourceUrl: 'https://github.com/kubernetes/dashboard',
    tags: ['visualization', 'configuration'],
  },
  {
    name: "keptn",
    shortDescription: 'Toolkit for cloud-native application lifecycle management',
    iconUrl: 'https://avatars.githubusercontent.com/u/46796476?s=150',
    websiteUrl: 'https://keptn.sh/',
    sourceUrl: 'https://github.com/keptn/lifecycle-toolkit',
    tags: ['delivery', 'logging'],
  },
  {
    name: "Cyclops",
    shortDescription: 'Developer friendly Kubernetes 👁️',
    iconUrl: 'https://cyclops-ui.com/img/logo.png',
    websiteUrl: 'https://cyclops-ui.com/',
    sourceUrl: 'https://github.com/cyclops-ui/cyclops',
    tags: ['visualization', 'configuration'],
  },
  {
    name: "K8sGPT",
    shortDescription: 'Automatic SRE Superpowers within your Kubernetes cluster',
    iconUrl: 'https://avatars.githubusercontent.com/u/128535266?s=150',
    websiteUrl: 'https://k8sgpt.ai/',
    sourceUrl: 'https://github.com/k8sgpt-ai/k8sgpt-operator',
    tags: ['planned', 'ai'],
  },
  {
    name: "Kube Prometheus",
    shortDescription: 'Use Prometheus to monitor Kubernetes and applications running on Kubernetes',
    iconUrl: 'https://avatars.githubusercontent.com/u/66682517?s=150',
    websiteUrl: 'https://prometheus-operator.dev/',
    sourceUrl: 'https://github.com/prometheus-operator/kube-prometheus',
    tags: ['planned', 'monitoring'],
  },
  {
    name: "Velero",
    shortDescription: 'Backup and migrate Kubernetes applications and their persistent volumes',
    iconUrl: 'https://github.com/glasskube/glasskube/assets/3041752/44efa073-fd5f-47b8-9ae5-95abb33a47e7',
    websiteUrl: 'https://velero.io/',
    sourceUrl: 'https://github.com/vmware-tanzu/velero',
    tags: ['planned', 'backup'],
  },
  {
    name: "CloudNativePG",
    shortDescription: 'CloudNativePG covers the full lifecycle of a PostgreSQL database cluster',
    iconUrl: 'https://avatars.githubusercontent.com/u/100373852?s=150',
    websiteUrl: 'https://cloudnative-pg.io/',
    sourceUrl: 'https://github.com/cloudnative-pg/cloudnative-pg',
    tags: ['planned', 'database'],
  },
  {
    name: "MariaDB Operator",
    shortDescription: 'Run and operate MariaDB in a cloud native way',
    iconUrl: 'https://avatars.githubusercontent.com/u/127887858?s=150',
    websiteUrl: 'https://github.com/mariadb-operator/mariadb-operator',
    sourceUrl: 'https://github.com/mariadb-operator/mariadb-operator',
    tags: ['planned', 'database'],
  },
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

export const Tags: { [type in TagType]: Tag } = {
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
};

export const TagList = Object.keys(Tags) as TagType[];

function sortUsers() {
  let result = Users;
  // Sort by site name
  result = sortBy(result, (user) => user.name.toLowerCase());
  // Sort by favorite tag, favorites first
  result = sortBy(result, (user) => user.tags.includes('planned'));
  return result;
  ;
}

export const sortedUsers = sortUsers();
