import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import {
  faCheckDouble,
  faHeartbeat,
  faKey,
  faPaintbrush,
  faSync,
  faWalkieTalkie,
} from '@fortawesome/free-solid-svg-icons';
import React from 'react';

type FeatureItem = {
  title: string;
  icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Monitoring',
    icon: faHeartbeat,
    description: (
      <>
        Monitoring ensures customers maintain visibility into the health status
        of their applications and infrastructure.
      </>
    ),
  },
  {
    title: 'Installation & Update Distribution',
    icon: faSync,
    description: (
      <>
        Installation & Update Distribution simplifies deployment processes and
        empowers customers with update control.
      </>
    ),
  },
  {
    title: 'Communication',
    icon: faWalkieTalkie,
    description: (
      <>
        Communication establishes the tools and methods needed to build reliable
        support channels. It bridges the gap between vendors and end customers,
        ensuring alignment and increasing the likelihood of mutual success.
      </>
    ),
  },
  {
    title: 'Application Design',
    icon: faPaintbrush,
    description: (
      <>
        Application Design forms the foundation of the entire structure. Without
        a solid “on-prem ready” application, the deployment process has few
        chances of succeeding.
      </>
    ),
  },
  {
    title: 'Deployment target validation',
    icon: faCheckDouble,
    description: (
      <>
        Deployment target validation is the step which if done correctly ensures
        that end-customers regardless of what environments they happen to call
        production will be compatible with the vendors software.
      </>
    ),
  },
  {
    title: 'Licensing',
    icon: faKey,
    description: (
      <>
        Licensing gives vendors control over how features are enabled or
        disabled for end-customers. Effective licensing mechanisms should
        reflect contract agreements in real-time, providing flexibility while
        safeguarding the vendor`s interests.
      </>
    ),
  },
];

function Feature({title, icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--6', 'margin-top--lg')}>
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

export default function SoftwareDistributionFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row margin-top--xl">
          <div className="col text--center">
            <Heading as="h2">
              Building blocks for Modern On-Prem Software distribution
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
