import React, {FC} from 'react';
import GitHubButton from 'react-github-btn';

interface CustomGitHubButtonProps {
  href: string;
}

const CustomGitHubButton: FC<CustomGitHubButtonProps> = ({href, ...props}) => (
  <GitHubButton
    href={href}
    data-color-scheme="no-preference: light; light: light; dark: light;"
    data-icon="octicon-star"
    data-size="large"
    data-show-count="true"
    aria-label={`Star ${href} on GitHub`}
    {...props}>
    Star
  </GitHubButton>
);

export default CustomGitHubButton;
