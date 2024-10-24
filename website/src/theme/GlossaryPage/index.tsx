import React from 'react';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import GlossaryItem from '@site/src/pages/glossary/_components/GlossaryItem';
import styles from './styles.module.css';

export default function GlossaryPage(): JSX.Element {
  const {siteConfig} = useDocusaurusContext();
  const glossaryItems = siteConfig.customFields.glossaryItems as { term: string; fileName: string }[];

  return (
    <main className="container margin-vert--lg">
      <h1>Glossary</h1>
      <p>Explore key terms and definitions related to Glasskube and Kubernetes.</p>
      <div className={styles.glossaryGrid}>
        {glossaryItems.map((item) => (
          <GlossaryItem key={item.term} term={item.term} />
        ))}
      </div>
    </main>
  );
}