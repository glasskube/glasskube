import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';
import React from 'react';
import {Player} from '@lottiefiles/react-lottie-player';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faStar} from '@fortawesome/free-regular-svg-icons';
import {faXTwitter} from '@fortawesome/free-brands-svg-icons';
import Typewriter from 'typewriter-effect';


function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();


  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <div className="row">
          <div className="col">
            <div className={styles.socialButtons}>
              <Link
                className="button button--secondary"
                to="https://github.com/glasskube/glasskube">
                <FontAwesomeIcon icon={faStar}/>&nbsp;Star
              </Link>
              <Link
                className="button button--secondary"
                to="https://x.com/intent/follow?screen_name=glasskube">
                <FontAwesomeIcon icon={faXTwitter}/>&nbsp;Follow
              </Link>
            </div>
          </div>
        </div>
        <div className="row row--algin-baseline">
          <div className="col padding-top--xl">
            <Heading as="h1" className="hero__title">
              <pre>
              <Typewriter
                onInit={(typewriter) => {
                  typewriter
                    .typeString('brew tap <span class="typewriter-command">glasskube/glasskube</span>')
                    .pauseFor(2500)
                    .deleteAll()
                    .typeString('brew install <span class="typewriter-command">glasskube</span>')
                    .pauseFor(2500)
                    .deleteAll()
                    .typeString('glasskube install ')
                    .typeString('<span class="typewriter-command">cert-manager</span>')
                    .pauseFor(2500)
                    .deleteChars('cert-manager'.length)
                    .pauseFor(2500)
                    .typeString('<span class="typewriter-command">ingress-nginx</span>')
                    .deleteChars('ingress-nginx'.length)
                    .pauseFor(2500)
                    .typeString('<span class="typewriter-command">kubernetes-dashboard</span>')
                    .deleteChars('kubernetes-dashboard'.length)
                    .pauseFor(2500)
                    .typeString('<span class="typewriter-command">[your-package]</span>')
                    .callFunction(() => {
                      console.log('All strings were deleted');
                    })
                    .start();
                }}
              />
                </pre>
              {/*glasskube install <span className={styles.pink}>cert-manager</span></pre>*/}
            </Heading>
            <p className="hero__subtitle">{siteConfig.tagline}</p>
            <div className={styles.buttons}>
              <Link
                className="button button--secondary button--lg"
                to="/docs/">
                🚀 Get started
              </Link>
              <Link
                className="button button--outline button--lg"
                to="https://discord.gg/SxH6KUCGH7">
                Join the community
              </Link>
            </div>
          </div>
          <div className="col">
            <Player
              autoplay
              loop
              src="/animations/home.json"
              style={{height: '480px'}}
            />
          </div>
        </div>

      </div>
    </header>
  );
}

export default function Home(): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title='Home'
      description={siteConfig.tagline}>
      <HomepageHeader/>
      <main>
        <HomepageFeatures/>
      </main>
    </Layout>
  );
}
