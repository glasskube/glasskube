import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {
  faGitAlt,
  faGithub,
  faSkyatlas,
} from '@fortawesome/free-brands-svg-icons';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import React from 'react';

type FeatureItem = {
  title: string;
  icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Enterprise ready',
    icon: faSkyatlas,
    description: (
      <>
        Manage the Kubernetes packages your core application depends on or
        distribute your Cloud Native application to your developers or customers
        with Glasskube.
      </>
    ),
  },
  {
    title: 'Advanced GitOps Integration',
    icon: faGitAlt,
    description: (
      <>
        Glasskube integrates into your GitOps workflow you already have in
        place. It integrates with Renovate and will provide resource level diffs
        right into your pull request.
      </>
    ),
  },
  {
    title: 'Open Source',
    icon: faGithub,
    description: (
      <>
        Glasskube is fully Open-Source, part of the CNCF landscape and is
        Apache-2.0 licensed. Developed by dozens of contributors from all over
        the world.
      </>
    ),
  },
];

function Feature({title, icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4', 'margin-top--lg')}>
      <div className={clsx('text--center', styles.iconBorder)}>
        <FontAwesomeIcon icon={icon} size="8x" className={styles.iconHeight} />
      </div>
      <div className="text--center padding-horiz--md margin-top--lg">
        <Heading as="h3" className="">
          {title}
        </Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function PackageManagerFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container margin-top--lg">
        <div className="row">
          <div className="col text--center">
            <Heading as="h2">
              Deploy, Configure and Update Kubernetes packages 20x faster than
              with Helm
            </Heading>
          </div>
        </div>
        <div className="row">
          {FeatureList.map((item, idx) => (
            <Feature key={idx} {...item} />
          ))}
        </div>
      </div>
    </section>
  );
}
