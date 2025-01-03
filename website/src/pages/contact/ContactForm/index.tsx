import Heading from '@theme/Heading';
import React from 'react';
import styles from './index.module.css';

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

interface ContactRequest {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  jobTitle: string;
  companyName: string;
  useCase: string;
}

class Form extends React.Component<unknown, {data: ContactRequest}> {
  hubSpotLoaded = false;

  constructor(props: unknown) {
    super(props);
    this.state = {data: {} as ContactRequest};

    this.handleSubmit = this.handleSubmit.bind(this);
  }

  loadHubSpot() {
    if (!this.hubSpotLoaded) {
      loadScript();
      this.hubSpotLoaded = true;
    }
  }

  async handleSubmit(event) {
    event.preventDefault();
    loadScript();

    await fetch('https://forms.glasskube.com/api/v1/contact', {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(this.state.data),
    })
      .then(() => {
        alert(
          'Thank you for reaching out. We will follow up as soon as possible.',
        );
      })
      .catch(err => {
        alert(err.message);
        console.error(err.message);
      });
  }

  render() {
    return (
      <form onSubmit={this.handleSubmit} className={styles.form}>
        <div className="row g-2 mt-1">
          <div className="col">
            <label htmlFor="firstName">First Name*</label>
            <input
              type="text"
              className={styles.input}
              id="firstName"
              name="firstName"
              placeholder="First name"
              value={this.state.data.firstName}
              onChange={e =>
                this.setState(state => {
                  state.data.firstName = e.target.value;
                })
              }
              required
            />
          </div>
          <div className="col">
            <label htmlFor="lastName">Last Name*</label>
            <input
              type="text"
              className={styles.input}
              id="lastName"
              name="lastName"
              placeholder="Last name"
              value={this.state.data.lastName}
              onChange={e =>
                this.setState(state => {
                  state.data.lastName = e.target.value;
                })
              }
              required
            />
          </div>
        </div>
        <div className="row g-2">
          <div className="col">
            <label htmlFor="email">Work Email*</label>
            <input
              type="email"
              className={styles.input}
              id="email"
              placeholder="Work Email"
              name="email"
              value={this.state.data.email}
              onChange={e => {
                this.setState(state => {
                  state.data.email = e.target.value;
                });
                this.loadHubSpot();
              }}
              required
            />
          </div>
          <div className="col">
            <label htmlFor="phone">Work Phone</label>
            <input
              type="tel"
              className={styles.input}
              id="phone"
              placeholder="Work Phone"
              value={this.state.data.phone}
              onChange={e =>
                this.setState(state => {
                  state.data.phone = e.target.value;
                })
              }
              name="phone"
            />
          </div>
        </div>
        <div className="row g-2">
          <div className="col">
            <label htmlFor="jobTitle">Job title*</label>
            <input
              type="text"
              className={styles.input}
              id="jobTitle"
              placeholder="Job Title"
              name="jobTitle"
              value={this.state.data.jobTitle}
              onChange={e =>
                this.setState(state => {
                  state.data.jobTitle = e.target.value;
                })
              }
              required
            />
          </div>
          <div className="col">
            <label htmlFor="companyName">Company Name*</label>
            <input
              type="text"
              className={styles.input}
              id="companyName"
              placeholder="Company Name"
              name="companyName"
              value={this.state.data.companyName}
              onChange={e =>
                this.setState(state => {
                  state.data.companyName = e.target.value;
                })
              }
              required
            />
          </div>
        </div>
        <div>
          <label htmlFor="useCase">Use case*</label>
          <textarea
            id="useCase"
            name="useCase"
            required
            placeholder=""
            value={this.state.data.email}
            onChange={e =>
              this.setState(state => {
                state.data.useCase = e.target.value;
              })
            }
            className={styles.input}
          />
        </div>
        <button className="button button--secondary button--lg" type="submit">
          Submit
        </button>
      </form>
    );
  }
}

export default function ContactForm() {
  return (
    <div className="container">
      <div className="row">
        <div className="col col--6 col--offset-3 margin-vert--lg">
          <div className="text--center">
            <Heading as={'h2'}>Contact form</Heading>
            <p>We will follow up as soon as possible.</p>
          </div>
          <Form />
        </div>
      </div>
    </div>
  );
}
