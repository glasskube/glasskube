import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faGitAlt} from '@fortawesome/free-brands-svg-icons';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import {faBoxes, faCodeBranch, faDisplay, faMagnifyingGlass, faSync} from '@fortawesome/free-solid-svg-icons';

type FeatureItem = {
  title: string;
  Icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'CLI and UI available',
    Icon: faDisplay,
    description: (
      <>
        CLI und UI as first class components. It doesn't matter if you prefer managing your
        packages via a CLI or UI - Glasskube supports both.
      </>
    ),
  },
  {
    title: 'Dependency aware',
    Icon: faCodeBranch,
    description: (
      <>
        Glasskube packages are dependency aware. If two packages require the same dependency,
        Glasskube makes sure it only gets installed onced.
      </>
    ),
  },
  {
    title: 'GitOps ready',
    Icon: faGitAlt,
    description: (
      <>
        All packages are stored in custom resources, which can easily managed with your favorit
        GitOps tool like ArgoCD or Flux.
      </>
    ),
  },
  {
    title: 'Automated updates',
    Icon: faSync,
    description: (
      <>
        Glasskube ensures your Kubernetes packages and apps are always up-to-date, minimizing the
        manual effort required for maintenance.
      </>
    ),
  },
  {
    title: 'Central package repository',
    Icon: faBoxes,
    description: (
      <>
        Keep track of all your packages in one central repository, which a planned feature for custom repositories.
        (planned)
      </>
    ),
  },
  {
    title: 'Cluster Scan',
    Icon: faMagnifyingGlass,
    description: (
      <>
        Introducing the Cluster Scan feature in a future version, which allows you to detect packages in your cluster,
        providing valuable insights for better management and upgrade paths.
      </>
    ),
  },
];

function Feature({title, Icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4', 'margin-top--lg')}>
      <div className="text--center">
        <FontAwesomeIcon icon={Icon} size="4x" className={styles.h64}/>
      </div>
      <div className="text--center padding-horiz--md margin-top--lg">
        <Heading as="h3" className="">{title}</Heading>
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
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
