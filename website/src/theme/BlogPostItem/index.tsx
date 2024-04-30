import React from 'react';
import BlogPostItem from '@theme-original/BlogPostItem';
import type BlogPostItemType from '@theme/BlogPostItem';
import type { WrapperProps } from '@docusaurus/types';
import { useBlogPost } from '@docusaurus/theme-common/internal';
import GiscusComponent from '@site/src/components/GiscusComponent';

type Props = WrapperProps<typeof BlogPostItemType>;

export default function BlogPostItemWrapper(props: Props) {
  const { metadata, isBlogPostPage } = useBlogPost();
  const { enableComments = true } = metadata.frontMatter;

  return (
    <>
      <BlogPostItem {...props} />
      {enableComments && isBlogPostPage && <GiscusComponent />}
    </>
  );
}
