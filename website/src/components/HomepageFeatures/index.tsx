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

type FeatureItem = {
  title: string;
  Icon: IconDefinition;
  description: JSX.Element;
};

const FeatureList: FeatureItem[] = [
  {
    title: 'Enterprise ready',
    Icon: faSkyatlas,
    description: (
      <>
        Manage the Kubernetes packages your core application depends on or
        distribute internal services charts to your developers with the
        Glasskube package manager.
      </>
    ),
  },
  {
    title: 'Advanced GitOps Integration',
    Icon: faGitAlt,
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
    Icon: faGithub,
    description: (
      <>
        Glasskube is fully Open-Source, part of the CNCF landscape and is
        Apache-2.0 licensed. Developed by dozens of contributors from all over
        the world.
      </>
    ),
  },
];

function Feature({title, Icon, description}: FeatureItem) {
  return (
    <div className={clsx('col col--4', 'margin-top--lg')}>
      <div className={clsx('text--center', styles.iconBorder)}>
        <FontAwesomeIcon icon={Icon} size="8x" className={styles.iconHeight} />
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
