import React from 'react';
import Link from '@docusaurus/Link';
import useBaseUrl from '@docusaurus/useBaseUrl';
import styles from './styles.module.css';
import {Content} from '@theme/BlogPostPage';
import clsx from 'clsx';

interface HomepageProps {
  homePageBlogMetadata: any;
  readonly recentPosts: readonly {readonly Preview: Content; metadata: any}[];
}

const LatestBlogPosts = ({recentPosts}: HomepageProps) => {
  return (
    <section className={styles.latestBlogPosts}>
      <div className="container margin-top--lg">
        <div className="row">
          <div className="col text--center">
            <h2>Check our blog</h2>
            <p>
              We are working on a blog post series covering topics from general
              open source to specific Kubernetes related topics.{' '}
              <Link to={useBaseUrl('/blog')}>Check it out!</Link>
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
              className={clsx(
                'button',
                'button--secondary',
                'button--lg',
                styles.moreBlogs,
              )}
              to={useBaseUrl('/blog')}>
              More blogs
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
};

export default LatestBlogPosts;
