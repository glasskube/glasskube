import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import React from 'react';
import {
  faArrowTrendUp,
  faClock,
  faLightbulb,
} from '@fortawesome/free-solid-svg-icons';

type FeatureItem = {
  title: string;
  icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Reduce onboarding time',
    icon: faClock,
    description: (
      <>
        Your customers receive a guided installation experience, even for
        diverse customer environments.
      </>
    ),
  },
  {
    title: 'Increase Update velocity',
    icon: faArrowTrendUp,
    description: (
      <>
        Glasskube Cloud notifies your customers about new updates and empowers
        them to update whenever they are ready.
      </>
    ),
  },
  {
    title: 'Simplify Troubleshooting',
    icon: faLightbulb,
    description: (
      <>
        Receive metrics, heart beat information and let your customers share
        additional support insights if needed.
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

export default function HomepageFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          <div className="col text--center margin-top--xl">
            <Heading as="h2">
              The Glasskube Cloud Software Distribution Platform is the central
              place to manage all your on-prem, vpc and air-gapped customers.
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
