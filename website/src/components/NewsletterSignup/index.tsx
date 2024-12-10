import Heading from '@theme/Heading';
import React, {ChangeEvent} from 'react';
import styles from '@site/src/components/Home/index.module.css';

function loadScript() {
  if (typeof window === 'undefined') {
    return null;
  }

  const elementId = 'hs-script';
  if (document.getElementById(elementId) === null) {
    const script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = 'https://js-eu1.hs-scripts.com/144345473.js';
    script.id = elementId;
    document.head.appendChild(script);
  }
}

class NewsletterForm extends React.Component<unknown, {value: string}> {
  constructor(props: unknown) {
    super(props);
    this.state = {value: ''};

    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(event: ChangeEvent) {
    if (!this.state.value) {
      loadScript();
    }
    this.setState({value: (event.target as HTMLInputElement).value});
  }

  async handleSubmit(event) {
    event.preventDefault();
    await fetch('https://cms.glasskube.eu/api/ezforms/submit', {
      method: 'POST',
      body: JSON.stringify({
        token: '',
        formName: 'newsletter',
        formData: {
          email: this.state.value,
        },
      }),
      headers: {
        'Content-type': 'application/json; charset=UTF-8',
      },
    })
      .then(response => response.json())
      .then(() => {
        alert(
          'Email successfully added to our newsletter: ' + this.state.value,
        );
      })
      .catch(err => {
        console.error(err.message);
      });
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit} className={styles.newsletterForm}>
        <input
          type="email"
          id="email"
          name="email"
          required
          placeholder="your-email@corp.com"
          value={this.state.value}
          onChange={this.handleChange}
          className={styles.emailInput}
        />
        <button className="button button--secondary button--lg" type="submit">
          Subscribe
        </button>
      </form>
    );
  }
}

export default function NewsletterSignup() {
  return (
    <div className="container text--center">
      <div className="row">
        <div className="col col--6 col--offset-3 margin-vert--lg">
          <div>
            <Heading as={'h2'}>Glasskube Newsletter</Heading>
            <p>Sign-Up to get the latest product updates and release notes!</p>
            <NewsletterForm />
          </div>
        </div>
      </div>
    </div>
  );
}
