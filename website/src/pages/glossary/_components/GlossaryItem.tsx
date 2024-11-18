import React from 'react';
import Link from '@docusaurus/Link';
import styles from './styles.module.css';

interface GlossaryItemProps {
  term: string;
  fileName: string;
}

export default function GlossaryItem({ term, fileName }: GlossaryItemProps): JSX.Element {
  return (
    <div className={styles.glossaryItem}>
      <div className={styles.glossaryItemInner}>
        <h3>{term}</h3>
        <p className={styles.description}>
          {/* Add a brief description here if available */}
        </p>
        <Link to={`/glossary/${fileName}`} className={styles.moreButton}>
          More
        </Link>
      </div>
    </div>
  );
}