import styles from './styles.module.css';
import CustomGitHubButton from '@site/src/components/CustomGitHubButton';


function Index() {
  return (
    <div className={styles.wrapper}>
      <div className={styles.center}>
        <CustomGitHubButton href='https://github.com/glasskube/glasskube'/>
      </div>

    </div>
  );
}

export default Index;
