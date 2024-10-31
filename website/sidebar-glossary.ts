import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  glossary: [
    {
      type: 'category',
      label: 'Glossary',
      link: {
        type: 'generated-index',
        title: 'Kubernetes Glossary',
        description:
          'Learn about common Kubernetes and cloud-native terminology',
        slug: '/categories/glossary',
      },
      items: [
        'devops',
        'helm',
        'helm-chart',
        'k8s-configuration-mgmt',
        'k8s-package-mgmt',
        'kubernetes',
        'kubernetes-operator',
      ], 
    },
  ],
};

export default sidebars;