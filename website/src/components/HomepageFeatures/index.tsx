import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import React from 'react';
import {
  faChartLine,
  faCheckDouble,
  faKey,
} from '@fortawesome/free-solid-svg-icons';

type FeatureItem = {
  title: string;
  icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Manage Licensing',
    icon: faKey,
    description: (
      <>
        Restrict and permit access to your software with built in license
        management and validation.
      </>
    ),
  },
  {
    title: 'Compatibility testing',
    icon: faCheckDouble,
    description: (
      <>
        Ensure compatibility across different environments like AWS, GCP and
        on-prem.
      </>
    ),
  },
  {
    title: 'Insights',
    icon: faChartLine,
    description: (
      <>
        Monitor uptime, usage, configuration, and issues of your enterprise
        deployments.
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
          <div className="col text--center">
            <Heading as="h2">
              Sell software to enterprises and deploy to private infrastructure,
              different cloud environments, and on-premises.
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
