import React, {FC} from 'react';
import Link from '@docusaurus/Link';

interface SignupForWaitlistButtonProps {
  additionalClassNames: string;
}

function posthogId() {
  if (window['posthog']) {
    return '/signup.html?id=' + window.posthog.get_distinct_id();
  }
  return '';
}

const SignupForWaitlistButton: FC<SignupForWaitlistButtonProps> = ({additionalClassNames}) => (
  <Link
    className={`button button--accent ${additionalClassNames}`}
    to={`https://glasskube.cloud${posthogId()}`}>
    Glasskube Cloud signup
  </Link>
);

export default SignupForWaitlistButton;
