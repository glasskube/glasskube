import styles from './styles.module.css';
import CustomGitHubButton from '@site/src/components/buttons/CustomGitHubButton';
import TalkToFoundersButton from '@site/src/components/buttons/TalkToFoundersButton';
import SignupForWaitlistButton from '@site/src/components/buttons/SignupForWaitlistButton';


function Index() {
  return (
    <>
      <div className={styles.wrapper}>
        <div className={styles.center}>
          <CustomGitHubButton href='https://github.com/glasskube/glasskube'/>
        </div>
      </div>
      <div className={styles.ctaWrapper}>
        <div className={styles.center}>
          <TalkToFoundersButton additionalClassNames={''}/>
          &nbsp;&nbsp;
          <SignupForWaitlistButton additionalClassNames={''}/>
        </div>
      </div>

    </>
  );
}

export default Index;
