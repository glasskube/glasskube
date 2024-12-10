import React from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';
import clsx from 'clsx';
import {HomepageProps} from '../Home';

const LatestBlogPosts = ({recentPosts}: HomepageProps) => {
  return (
    <section className={styles.latestBlogPosts}>
      <div className="container margin-top--lg">
        <div className="row">
          <div className="col text--center">
            <h2>Recent blog posts</h2>
            <p>
              We regularly write{' '}
              <Link to={useBaseUrl('/blog/')}>blog posts</Link> and publish{' '}
              <Link to={useBaseUrl('/guides/')}>guides</Link>. Find our latest
              blog posts below.
            </p>
          </div>
        </div>
        <div className={styles.blogList}>
          {recentPosts.map((post, index) => (
            <div key={index} className={clsx('card', styles.blogCard)}>
              <div className={clsx('card__image', styles.blogCardImage)}>
                <img src={post.metadata.frontMatter.image} />
              </div>
              <div className={clsx(styles.cardContent, 'card__body')}>
                <img
                  src={post.Preview.metadata.authors[0].imageURL}
                  alt={post.Preview.metadata.authors[0].name}
                  className={styles.authorImage}
                />
                <div className={styles.cardDesc}>
                  <h3>
                    <Link to={post.Preview.metadata.permalink}>
                      {post.metadata.title}
                    </Link>
                  </h3>
                  <div>{post.metadata.description}</div>
                </div>
              </div>
            </div>
          ))}
        </div>
        <div className="row">
          <div className="col text--center">
            <Link
              className={clsx('button', 'button--secondary', styles.moreBlogs)}
              to={useBaseUrl('/blog')}>
              More posts
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
};

export default LatestBlogPosts;
