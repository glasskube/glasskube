import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  guides: [
    {
      type: 'category',
      label: 'Package Installation Guides',
      link: {
        type: 'generated-index',
        title: 'Package Installation Guides',
        description:
          '⚠️ Learn how to install your favorite Kubernetes add-ons using the Glasskube package manager ⚠️',
        slug: '/',
      },
      items: [
        'cert-manager',
        'ingress-nginx',
        'rabbitmq',
        'quickwit',
        'kube-prom-stack',
        'cloudnativepg'
      ],
    },
    {
      type: 'category',
      label: 'Contributor Guides',
      link: {
        type: 'generated-index',
        title: 'Contributor Guides',
        description:
          '⚠️ Contributors are what make open source great, here is where we share some resources to support you. ⚠️',
        slug: '/contributors',
      },
      items: ['git-guide', 'kubectl-guide', 'package-creation'],
    },
  ],
};

export default sidebars;
