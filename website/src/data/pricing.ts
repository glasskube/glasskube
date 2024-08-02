/* eslint-disable global-require */

export const PricingFAQs: PricingFAQ[] = [
  {
    id: 'free-trial',
    question: 'Do you offer free trials?',
    answer:
      'Please <a href="https://cal.glasskube.eu/team/founder/enterprise" target="_blank">contact us</a> regarding the enterprise mode for an individual offer.',
  },
  {
    id: 'free-tier',
    question: 'How are 5 nodes defined?',
    answer:
      'Only worker nodes are counting towards your free node quote. Control plane nodes are not counted in.',
  },
  {
    id: 'invoice',
    question: 'Do you also support pay by invoice?',
    answer:
      'Yes, individual offers from the enterprise plan can also be paid by invoice.',
  },
  {
    id: 'currency',
    question: 'Do you also support other currencies than US Dollar?',
    answer:
      'Yes, individual offers from the enterprise plan can be created in different currencies.',
  },
  {
    id: 'gitops',
    question: 'Does Glasskube integrate with existing GitOps solutions?',
    answer:
      'Yes, Glasskube integrates into your favorite GitOps solution, by creating pull requests for package upgrades.',
  },
  {
    id: 'paas',
    question: 'Is Glasskube a PaaS solution?',
    answer:
      'Glasskube helps you package and deploy Cloud Native Applications into your own Kubernetes cluster. We are not a PaaS solution built on top of Kubernetes, but integrate into your existing workflows.',
  },
  {
    id: 'kubernetes',
    question: 'Do I need to manage the Kubernetes cluster myself?',
    answer:
      'Generally yes, please <a href="https://cal.glasskube.eu/team/founder/enterprise" target="_blank">contact us</a> regarding your individual challenges.',
  },
];

export type PricingFAQ = {
  id: string;
  question: string;
  answer: string;
};
