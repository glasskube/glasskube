import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';
import Heading from '@theme/Heading';

import styles from './index.module.css';
import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faStar } from '@fortawesome/free-regular-svg-icons';
import { faXTwitter } from '@fortawesome/free-brands-svg-icons';
import Typewriter from 'typewriter-effect';
import BrowserOnly from '@docusaurus/BrowserOnly';


function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext();

  return (
    <header className={clsx('hero hero--primary', styles.heroBanner)}>
      <div className="container">
        <div className="row">
          <div className="col">
            <div className={styles.socialButtons}>
              <Link
                className={clsx('button', 'button--secondary', styles.socialButtonsLink)}
                to="https://github.com/glasskube/glasskube">
                <FontAwesomeIcon icon={faStar} />&nbsp;Star
              </Link>
              <Link
                className={clsx('button', 'button--secondary', styles.socialButtonsLink)}
                to="https://x.com/intent/follow?screen_name=glasskube">
                <FontAwesomeIcon icon={faXTwitter} />&nbsp;Follow
              </Link>
            </div>
          </div>
        </div>
        <div className="row row--no-gutters">
          <div className={clsx('col', styles.heroCol)}>
            <Heading as="h1" className="hero__title">
              <pre>
                <Typewriter
                  onInit={(typewriter) => {
                    typewriter
                      .changeDeleteSpeed(25)
                      .changeDelay(75)
                      .typeString('brew tap <span class="typewriter-command">glasskube/glasskube</span>')
                      .pauseFor(1500)
                      .deleteAll(25)
                      .typeString('brew install <span class="typewriter-command">glasskube</span>')
                      .pauseFor(1500)
                      .deleteAll(25)
                      .typeString('glasskube install ')
                      .typeString('<span class="typewriter-command">cert-manager</span>')
                      .pauseFor(1500)
                      .deleteChars('cert-manager'.length)
                      .typeString('<span class="typewriter-command">ingress-nginx</span>')
                      .pauseFor(1500)
                      .deleteChars('ingress-nginx'.length)
                      .typeString('<span class="typewriter-command">kubernetes-dashboard</span>')
                      .deleteChars('kubernetes-dashboard'.length)
                      .typeString('<span class="typewriter-command">[your-package]</span>')
                      .start();
                  }}
                />
              </pre>
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
            <div className={styles.lottiePlayerWrapper}>
              <BrowserOnly fallback={<div className={styles.lottieFallback}>Loading...</div>}>
                {() => {
                  const Player = require('@lottiefiles/react-lottie-player').Player;
                  return <Player
                    autoplay
                    loop
                    src="/animations/home.json"
                    style={{ height: '480px' }}
                  />
                }}
              </BrowserOnly>
            </div>
          </div>
        </div>

      </div>
    </header>
  );
}

function HomepageVideo() {
  const { siteConfig } = useDocusaurusContext();

  return (
    <div className={clsx('container-fluid', 'text--center', styles.backgroundSecondary)}>
      <div className="container text--center">
        <div className="row">
          <div className="col col--6 col--offset-3 margin-vert--lg">
            <div>
              <Heading as={'h2'} className={styles.colorWhite}>
                Learn how to install cert-manager in less than 2 minutes using glasskube
              </Heading>
              <a href="https://asciinema.org/a/k19wlsoX5Mr3raY6ro13imyNo" target="_blank"><img src="https://asciinema.org/a/k19wlsoX5Mr3raY6ro13imyNo.svg" /></a>
            </div>
          </div>
        </div>
      </div>
    </div>

  );
}

function HomepageNewsletter() {
  const { siteConfig } = useDocusaurusContext();

  return (
    <div className="container text--center">
      <div className="row">
        <div className="col col--6 col--offset-3 margin-vert--lg">
          <div>
            <Heading as={'h2'}>
              Glasskube Newsletter
            </Heading>
            <p>Sign-Up to get the latest product updates and release notes!</p>
            <NewsletterForm />
          </div>
        </div>
      </div>
    </div>

  );
}

class NewsletterForm extends React.Component<any, { value: string }> {
  constructor(props) {
    super(props);
    this.state = { value: '' };

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event) {
    this.setState({ value: event.target.value });
  }

  async handleSubmit(event) {
    event.preventDefault()
    await fetch('https://cms.glasskube.eu/api/ezforms/submit', {
      method: 'POST',
      body: JSON.stringify({
        token: '',
        formName: 'newsletter',
        formData: {
          email: this.state.value,
        }
      }),
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
      },
    })
      .then((response) => response.json())
      .then((data) => {
        alert('Email successfully added to our newsletter: ' + this.state.value);
      })
      .catch((err) => {
        console.error(err.message);
      });
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit} className={styles.newsletterForm}>
        <input type="email" id="email" name="email"
          placeholder="your-email@corp.com"
          value={this.state.value} onChange={this.handleChange}
          className={styles.emailInput} />
        <button
          className="button button--secondary button--lg"
          type="submit">
          Subscribe
        </button>
      </form>
    );
  }
}

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title='Home'
      description={siteConfig.tagline}>
      <HomepageHeader />
      <main>
        <HomepageFeatures />
        <HomepageVideo />
        <HomepageNewsletter />
      </main>
    </Layout>
  );
}
