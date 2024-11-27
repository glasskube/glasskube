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
        Glasskube is the easiest way to distribute and manage your software to
        your enterprise customers and edge locations.
      </p>
      <div className={styles.buttons}>
        <SignupForWaitlistButton additionalClassNames="" />
        <DemoButton additionalClassNames="" />
      </div>
    </div>
  );
}
