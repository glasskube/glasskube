import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faGitAlt} from '@fortawesome/free-brands-svg-icons';
import {IconDefinition} from '@fortawesome/free-regular-svg-icons';
import {faCodeBranch, faDisplay} from '@fortawesome/free-solid-svg-icons'; // Import the FontAwesomeIcon component.

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
];

function Feature({title, Icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4', 'margin-top--lg')}>
      <div className="text--center">
        <FontAwesomeIcon icon={Icon} size="4x" className={styles.h64} />
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
