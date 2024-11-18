export const PricingFAQs: PricingFAQ[] = [
  {
    id: 'free-trial',
    question: 'Do you offer free trials?',
    answer:
      'Please <a href="https://cal.glasskube.eu/team/founder/enterprise" target="_blank">contact us</a> regarding the enterprise mode for an individual offer.',
  },
  {
    id: 'free-tier',
    question: 'How are 5 customer deployments defined?',
    answer:
      'Only deployments with different license keys are counting towards your included license quota.',
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
      'Yes, Glasskube integrates with popular GitOps solution, by creating pull requests for package upgrades.',
  },
];

export type PricingFAQ = {
  id: string;
  question: string;
  answer: string;
};
