import React, {FC} from 'react';
import Link from '@docusaurus/Link';
import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment';

interface SignupForWaitlistButtonProps {
  additionalClassNames: string;
}

function posthogId() {
  if (ExecutionEnvironment.canUseDOM && window['posthog']) {
    try {
      return '/signup.html?id=' + window.posthog.get_distinct_id();
    } catch (e) {
      // no id
    }
  }
  return '';
}

const SignupForWaitlistButton: FC<SignupForWaitlistButtonProps> = ({
  additionalClassNames,
}) => (
  <Link
    className={`glasskube-telemetry-waitlist button button--outline ${additionalClassNames}`}
    to={`https://glasskube.cloud${posthogId()}`}>
    Glasskube Cloud waitlist
  </Link>
);

export default SignupForWaitlistButton;
