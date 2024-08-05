import Admonition from '@theme/Admonition';
import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import {
  faChartLine,
  faCheckDouble,
  faCodeCompare,
  faComments,
  faCubes,
  faGears,
} from '@fortawesome/free-solid-svg-icons';
import Link from '@docusaurus/Link';
import React from 'react';

type FeatureItem = {
  title: string;
  icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Insights and Reporting',
    icon: faChartLine,
    description: (
      <>
        Gain valuable insights with our comprehensive reporting dashboard.
        Monitor uptime, usage, and configuration of all private packages. Stay
        informed with notifications and maintain visibility into your software`s
        performance. Interested in implementing usage-based pricing? We can help
        with that too.
      </>
    ),
  },
  {
    title: 'Standardize Releases',
    icon: faCodeCompare,
    description: (
      <>
        Standardize your release process with an automated approach. Efficiently
        distribute new versions to various customer segments, ensuring each
        receives the correct software with the appropriate features. Manage
        releases effortlessly and maintain consistency across deployments.
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
    title: 'Resolve Technical Issues with Playbooks',
    icon: faComments,
    description: (
      <>
        Easily troubleshoot with shared support bundles between you and your
        customers. Integrated playbooks can trigger automated processes like
        database major version upgrades.
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
        to a specific version or version range, such as for Cert Manager or
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
      <div className="container margin-top--md">
        <div className="row">
          <div className="col col--6">
            <div className="card shadow--md">
              <div className="card__header text--center">
                <Heading as="h2">Open Source Packages</Heading>
              </div>
              <div className="card__body">
                <p>
                  Access a comprehensive library of Open Source packages on
                  Glasskube Hub. These packages are publicly available and can
                  be easily installed with the{' '}
                  <Link to="/products/package-manager">
                    Glasskube Package Manager
                  </Link>
                  . Eliminate the need for directly installing Helm charts and
                  manually applying manifests. With each update tested in our
                  public CI/CD pipelines, you can deploy with confidence and
                  upgrade seamlessly.
                  <br />
                  <br />
                  Missing a package?{' '}
                  <Link to="https://github.com/glasskube/glasskube/discussions/90">
                    Submit it here.
                  </Link>
                </p>
              </div>
            </div>
          </div>
          <div className="col col--6">
            <div className="card shadow--md">
              <div className="card__header text--center">
                <Heading as="h2">Private Packages</Heading>
              </div>
              <div className="card__body">
                <Admonition type="tip" icon="💡" title="Limited Offer">
                  The first 5 private repositories include free end-to-end
                  setup!
                </Admonition>
                <p>
                  Glasskube Hub offers secure private repositories to deploy
                  your application into your customers` infrastructure. Together
                  we create a Glasskube package for your application, host it in
                  a private secure repository, and configure it according to
                  your customers` Kubernetes-based infrastructure requirements.
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="row margin-top--xl">
          <div className="col text--center">
            <Heading as="h2">
              Glasskube Hub Enterprise for Private Packages
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
