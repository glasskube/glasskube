import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import {
  faChartLine,
  faCheckDouble,
  faCloud,
  faCubes,
  faGears,
  faKey,
} from '@fortawesome/free-solid-svg-icons';
import React from 'react';

type FeatureItem = {
  title: string;
  icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Integrate with Glasskube Cloud',
    icon: faCloud,
    description: (
      <>
        Integrate your private package repositories with Glasskube Cloud, to
        distribute your software to Kubernetes clusters.
      </>
    ),
  },
  {
    title: 'Insights and Reporting',
    icon: faChartLine,
    description: (
      <>
        Gain valuable insights with our comprehensive reporting dashboard.
        Monitor usage, and configuration of all private packages. Stay informed
        with notifications and maintain visibility into your software`s
        performance. Interested in implementing usage-based pricing? We can help
        with that too.
      </>
    ),
  },
  {
    title: 'License Management',
    icon: faKey,
    description: (
      <>
        Control access to your software with built-in license management and
        validation. Easily configure, validate, and monitor licenses to ensure
        only authorized users can access your software via a private repository
        hosted by Glasskube.
      </>
    ),
  },
  {
    title: 'Comprehensive Compatibility Testing',
    icon: faCheckDouble,
    description: (
      <>
        Ensure the reliability of your private packages with extensive testing
        across various cloud environments and Kubernetes versions. Whether your
        customer uses AWS, GCP, or on-premises, we ensure your application
        works!
      </>
    ),
  },
  {
    title: 'Installation Processes and configuration with Profiles',
    icon: faGears,
    description: (
      <>
        Meet your customers` diverse installation needs with Glasskube Hub.
        Offer different configuration options to cater to every level of
        Kubernetes expertise within your customer base. Ship private packages
        with different configurations and values for the same application, such
        as one package designed for 500+ users and another for 10k+ users.
      </>
    ),
  },
  {
    title: 'Build on Top of Our Comprehensive Open Source Package Library',
    icon: faCubes,
    description: (
      <>
        Don`t reinvent the wheel. All open-source dependencies of your
        application can be sourced from the Glasskube Hub`s Open Source package
        library. These packages are fully tested, and your application can point
        to a specific version or version range, such as for cert-manager or
        RabbitMQ.
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

export default function HubFeatures(): JSX.Element {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row margin-top--xl">
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
