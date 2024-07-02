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
        slug: '/categories/',
      },
      items: ['cert-manager', 'ingress-nginx', 'rabbitmq', 'quickwit'],
    },
    {
      type: 'category',
      label: 'Contributor Guides',
      link: {
        type: 'generated-index',
        title: 'Contributor Guides',
        description:
          '⚠️ Contributors are what make open source great, here is where we share some resources to support you. ⚠️',
        slug: '/categories/contributors',
      },
      items: ['git-guide', 'kubectl-guide'],
    },
  ],
};

export default sidebars;
