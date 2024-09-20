import Link from '@docusaurus/Link';
import DemoButton from '@site/src/components/buttons/DemoButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';
import React from 'react';
import styles from './styles.module.css';
import ContactSalesButton from '@site/src/components/buttons/ContactSalesButton';

function Pricing() {
  return (
    <section>
      <div className="container">
        <div className="row">
          <div className="col col--4">
            <div className="card margin-top--lg shadow--md">
              <div className="card__header text--center">
                <h3>Open Source Kubernetes Package Manager</h3>
                <h4>FREE</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Ideal for individuals and companies building and managing
                  their own pipelines and workflows utilizing the Glasskube
                  Package Manager.
                </p>
                <ul className={styles.pricingItemList}>
                  <li>Apache 2.0 licensed</li>
                  <li>Package installation CLI and UI</li>
                  <li>Basic GitOps Integration with Renovate</li>
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
            <div className="card shadow--md">
              <div className="card__header text--center">
                <h3>Application Delivery Platform for Kubernetes</h3>
                First 5 customer deployments included
                <br />
                <h4>Starting at $18k Platform fee per year</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Streamline how your software gets distributed into customer
                  managed Kubernetes environments in private infrastructure,
                  different cloud environments, and on-premises (even
                  air-gapped).
                </p>
                <ul className={styles.pricingItemList}>
                  <li>Integrated license management</li>
                  <li>Insights and Reporting</li>
                  <li>Comprehensive Compatibility Testing</li>
                  <li>Resolve Technical Issues with Playbooks</li>
                  <li>
                    Installation Processes with configuration options via GUI or
                    CLI
                  </li>
                  <li>BYOC</li>
                </ul>
              </div>
              <div className="card__footer">
                <SignupForWaitlistButton
                  additionalClassNames={'button--lg button--block'}
                />
                <div className="margin-top--md" />
                <DemoButton additionalClassNames={'button--lg button--block'} />
              </div>
            </div>
          </div>
          <div className="col">
            <div className="card col--4 margin-top--lg shadow--md">
              <div className="card__header text--center">
                <h3>Enterprise hardened Kubernetes Packages</h3>
                <h4>Custom quote</h4>
              </div>
              <hr />
              <div className="card__body">
                <p>
                  Tailor-made solutions for the highest level of cloud security
                  and reliability and personal support.
                </p>
                <p>Glasskube offers:</p>
                <ul className={styles.pricingItemList}>
                  <li>Free DevOps assessment call</li>
                  <li>SLAs and personal onboarding</li>
                  <li>Dedicated success manager</li>
                </ul>
              </div>
              <div className="card__footer">
                <ContactSalesButton
                  additionalClassNames={'button--lg button--block'}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

export default React.memo(Pricing);
