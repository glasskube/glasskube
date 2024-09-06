import Link from '@docusaurus/Link';
import React from 'react';

export default function BlogSidebar(): JSX.Element {
  return (
    <div>
        <div style={{fontSize:"28px"}}>
        Recent posts
        </div>
        <div>
        <Link
        style={{marginTop:'20px'}}
        className={`button button--outline`}
        to="/blog">
        View Recent posts
      </Link>
        </div>
    </div>
  );
}
