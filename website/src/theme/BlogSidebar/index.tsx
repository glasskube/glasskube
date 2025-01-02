import React from 'react';
import styles from './styles.module.css';
import clsx from 'clsx';
import DemoButton from '@site/src/components/buttons/DemoButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';

export default function BlogSidebar(): JSX.Element {
  return (
    <div className={clsx('card__header', styles.sidebar)}>
      <h3>About Glasskube</h3>
      <p>
        We are the creators of the Glasskube Open-Source Package Manager for
        Kubernetes and offering a comprehensive yet easy to use software
        distribution platform.
      </p>
      <div className={styles.buttons}>
        <DemoButton additionalClassNames="" />
        <SignupForWaitlistButton additionalClassNames="" />
      </div>
    </div>
  );
}
