import React, {FC} from 'react';
import Link from '@docusaurus/Link';

interface ContactSalesButtonProps {
  additionalClassNames: string;
}

const ContactSalesButton: FC<ContactSalesButtonProps> = ({
  additionalClassNames,
}) => (
  <Link
    className={`glasskube-telemetry-sales button button--info ${additionalClassNames}`}
    to="https://cal.glasskube.eu/team/founder/enterprise">
    Contact sales
  </Link>
);

export default ContactSalesButton;
