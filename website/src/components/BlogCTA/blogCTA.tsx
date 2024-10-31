import React from 'react';
import DemoButton from '../buttons/DemoButton';
import TalkToFoundersButton from '../buttons/TalkToFoundersButton';
import styles from './styles.module.css';
import clsx from 'clsx';

export default function BlogCTA(): JSX.Element {
  return (
    <div className={clsx('margin-top--xl margin-bottom--lg', styles.ctaWrapper)}>
      <div className="container">
        <div className="row">
          <div className="col col--10 col--offset-1">
            <div className={styles.ctaContent}>
              <h2>Issues with private infra and On-prem deployments?</h2>
              <p className="margin-bottom--lg">
                Book a call with us to discuss your Kubernetes deployment and scaling needs.
              </p>
              <div className={styles.buttonContainer}>
                <DemoButton additionalClassNames="margin-right--md" />
                <TalkToFoundersButton additionalClassNames="" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}