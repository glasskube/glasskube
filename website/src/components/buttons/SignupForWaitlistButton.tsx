import React, {FC} from 'react';
import Link from '@docusaurus/Link';

interface SignupForWaitlistButtonProps {
  additionalClassNames: string;
}

const SignupForWaitlistButton: FC<SignupForWaitlistButtonProps> = ({additionalClassNames}) => (
  <Link
    className={`button button--accent ${additionalClassNames}`}
    to="https://glasskube.cloud/">
    Signup for the wait list
  </Link>
);

export default SignupForWaitlistButton;
