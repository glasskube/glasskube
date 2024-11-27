import React from 'react';
import GlossaryItem from '@site/src/pages/glossary/_components/GlossaryItem';
import styles from './styles.module.css';
import sidebars from '@site/sidebar-glossary';

export default function GlossaryPage(): JSX.Element {
  const glossaryItems = sidebars.glossary[0].items as string[];

  return (
    <main className="container margin-vert--lg">
      <h1>Glossary</h1>
      <p>
        Explore key terms and definitions related to Glasskube and Kubernetes.
      </p>
      <div className={styles.glossaryGrid}>
        {glossaryItems.map(fileName => (
          <GlossaryItem
            key={fileName}
            term={fileName
              .split('-')
              .map(word => word.charAt(0).toUpperCase() + word.slice(1))
              .join(' ')}
            fileName={fileName}
          />
        ))}
      </div>
    </main>
  );
}
