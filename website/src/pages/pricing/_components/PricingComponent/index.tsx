import Link from '@docusaurus/Link';
import DemoButton from '@site/src/components/buttons/DemoButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';
import React from 'react';
import styles from './styles.module.css';

function Pricing() {
  return (
    <section>
      <div className="container">
        <div className="row">
          <div className="col col--4">
            <div className="card margin-top--lg">
              <div className="card__header text--center">
                <h3>Totally free</h3>
                <h4>$0</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Ideal for individuals and companies building and managing their own pipelines and workflows utilizing
                  the Glasskube Package Manager.
                </p>
                <ul className={styles.pricingItemList}>
                  <li>Apache 2.0 licensed</li>
                  <li>Package installation CLI and UI</li>
                  <li>GitHub & Discord community support</li>
                  <li>
                    Use our public package repository or host private packages
                    yourself
                  </li>
                </ul>
              </div>
              <div className="card__footer">
                <Link
                  className="button button--secondary button--block button--lg"
                  to="https://github.com/glasskube/glasskube/">
                  Start on GitHub
                </Link>
              </div>
            </div>
          </div>
          <div className="col col--4">
            <div className="card">
              <div className="card__header text--center">
                <h3>Best value champion</h3>
                Free for clusters with less than 5 nodes
                <br />
                <h4>Starting at $12k Platform fee per year</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Accelerate your company&apos;s DevOps workflows with private
                  hosted packages and improved cloud security.
                </p>
                <ul className={styles.pricingItemList}>
                  <li>
                    GitOps Integration with pull request enhancement and change
                    preview on a manifest level
                  </li>
                  <li>Backup and restore your installed packages</li>
                  <li>
                    Private hosted packages for secure package distribution
                  </li>
                  <li>
                    Security alerts and update notifications in case of
                    vulnerabilities in your Kubernetes packages for improved
                    software supply chain security
                  </li>
                  <li>Secure cloud access with RBAC for team members</li>
                </ul>
              </div>
              <div className="card__footer">
                <SignupForWaitlistButton
                  additionalClassNames={'button--lg button--block'}
                />
              </div>
            </div>
          </div>
          <div className="col">
            <div className="card col--4 margin-top--lg">
              <div className="card__header text--center">
                <h3>Enterprise mode</h3>
                <h4>Custom quote</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Tailor-made solutions for the highest level of cloud security
                  and reliability and personal support.
                </p>
                <p>All Glasskube Cloud features and:</p>
                <ul className={styles.pricingItemList}>
                  <li>Free DevOps assessment call</li>
                  <li>SLAs and personal onboarding</li>
                  <li>Dedicated success manager</li>
                </ul>
              </div>
              <div className="card__footer">
                <DemoButton additionalClassNames={'button--lg button--block'} />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default React.memo(Pricing);
