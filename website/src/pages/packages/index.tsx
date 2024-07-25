import {useEffect, useMemo, useState} from 'react';
import clsx from 'clsx';
import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment';
import {translate} from '@docusaurus/Translate';
import {useHistory, useLocation} from '@docusaurus/router';

import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import {type Package, sortedUsers, TagList, Tags, type TagType} from '@site/src/data/packages';
import Heading from '@theme/Heading';
import PackageTagSelect, {readSearchTags,} from './_components/PackageTagSelect';
import PackageFilterToggle, {type Operator, readOperator} from './_components/PackageFilterToggle';
import PackageCard from './_components/PackageCard';

import styles from './styles.module.css';

const TITLE = translate({message: 'Glasskube package repository'});
const DESCRIPTION = translate({
  message: 'List of packages that are or will be installable via the glasskube package manager',
});
const SUBMIT_URL = 'https://github.com/glasskube/glasskube/discussions/90';

type UserState = {
  scrollTopPosition: number;
  focusedElementId: string | undefined;
};

function restoreUserState(userState: UserState | null) {
  const {scrollTopPosition, focusedElementId} = userState ?? {
    scrollTopPosition: 0,
    focusedElementId: undefined,
  };
  document.getElementById(focusedElementId)?.focus();
  window.scrollTo({top: scrollTopPosition});
}

export function prepareUserState(): UserState | undefined {
  if (ExecutionEnvironment.canUseDOM) {
    return {
      scrollTopPosition: window.scrollY,
      focusedElementId: document.activeElement?.id,
    };
  }

  return undefined;
}

const SearchNameQueryKey = 'name';

function readSearchName(search: string) {
  return new URLSearchParams(search).get(SearchNameQueryKey);
}

function filterUsers(
  users: Package[],
  selectedTags: TagType[],
  operator: Operator,
  searchName: string | null,
) {
  if (searchName) {
    // eslint-disable-next-line no-param-reassign
    users = users.filter((user) =>
      user.name.toLowerCase().includes(searchName.toLowerCase()),
    );
  }
  if (selectedTags.length === 0) {
    return users;
  }
  return users.filter((user) => {
    if (user.tags.length === 0) {
      return false;
    }
    if (operator === 'AND') {
      return selectedTags.every((tag) => user.tags.includes(tag));
    }
    return selectedTags.some((tag) => user.tags.includes(tag));
  });
}

function useFilteredUsers() {
  const location = useLocation<UserState>();
  const [operator, setOperator] = useState<Operator>('OR');
  // On SSR / first mount (hydration) no tag is selected
  const [selectedTags, setSelectedTags] = useState<TagType[]>([]);
  const [searchName, setSearchName] = useState<string | null>(null);
  // Sync tags from QS to state (delayed on purpose to avoid SSR/Client
  // hydration mismatch)
  useEffect(() => {
    setSelectedTags(readSearchTags(location.search));
    setOperator(readOperator(location.search));
    setSearchName(readSearchName(location.search));
    restoreUserState(location.state);
  }, [location]);

  return useMemo(
    () => filterUsers(sortedUsers, selectedTags, operator, searchName),
    [selectedTags, operator, searchName],
  );
}

function PackagesHeader() {
  return (
    <section className="margin-top--lg margin-bottom--lg text--center">
      <Heading as="h1">{TITLE}</Heading>
      <p>{DESCRIPTION}</p>
      <Link className="button button--primary" to={SUBMIT_URL}>
        🙏 Please submit your package 📦
      </Link>
    </section>
  );
}

function PackagesFilters() {
  const filteredUsers = useFilteredUsers();
  return (
    <section className="container margin-top--l margin-bottom--lg">
      <div className={clsx('margin-bottom--sm', styles.filterCheckbox)}>
        <div>
          <Heading as="h2">
            Filters
          </Heading>
          <span>{filteredUsers.length} packages</span>
        </div>
        <PackageFilterToggle/>
      </div>
      <ul className={clsx('clean-list', styles.checkboxList)}>
        {TagList.map((tag, i) => {
          const {label, color} = Tags[tag];
          const id = `packages_checkbox_id_${tag}`;

          return (
            <li key={i} className={styles.checkboxListItem}>

              <PackageTagSelect
                tag={tag}
                id={id}
                label={label}
                icon={
                  <span
                    style={{
                      backgroundColor: color,
                      width: 10,
                      height: 10,
                      borderRadius: '50%',
                      marginLeft: 8,
                    }}
                  />
                }
              />
            </li>
          );
        })}
      </ul>
    </section>
  );
}

function SearchBar() {
  const history = useHistory();
  const location = useLocation();
  const [value, setValue] = useState<string | null>(null);
  useEffect(() => {
    setValue(readSearchName(location.search));
  }, [location]);
  return (
    <div className={styles.searchContainer}>
      <input
        id="searchbar"
        placeholder="Search for packages ..."
        value={value ?? undefined}
        onInput={(e) => {
          setValue(e.currentTarget.value);
          const newSearch = new URLSearchParams(location.search);
          newSearch.delete(SearchNameQueryKey);
          if (e.currentTarget.value) {
            newSearch.set(SearchNameQueryKey, e.currentTarget.value);
          }
          history.push({
            ...location,
            search: newSearch.toString(),
            state: prepareUserState(),
          });
          setTimeout(() => {
            document.getElementById('searchbar')?.focus();
          }, 0);
        }}
      />
    </div>
  );
}

function Packages() {
  const filteredUsers = useFilteredUsers();

  if (filteredUsers.length === 0) {
    return (
      <section className="margin-top--lg margin-bottom--xl">
        <div className="container padding-vert--md text--center">
          <Heading as="h2">
            No result
          </Heading>
        </div>
      </section>
    );
  }

  return (
    <section className="margin-top--lg margin-bottom--xl">
      {filteredUsers.length === sortedUsers.length ? (
        <>
          <div className={styles.supportedPackages}>
            <div className="container">
              <div
                className={clsx(
                  'margin-bottom--md',
                  styles.supportedPackagesHeader,
                )}>
                <Heading as="h2">
                  Already supported
                </Heading>
              </div>
              <ul
                className={clsx('clean-list', styles.packageList,)}>
                {favoriteUsers.map((user) => (
                  <PackageCard key={user.name} user={user}/>
                ))}
              </ul>
            </div>
          </div>
          <div className="container margin-top--lg">
            <Heading as="h2">
              Planned
            </Heading>
            <ul className={clsx('clean-list', styles.packageList)}>
              {otherUsers.map((user) => (
                <PackageCard key={user.name} user={user}/>
              ))}
            </ul>
          </div>
        </>
      ) : (
        <div className="container">
          <div
            className={clsx('margin-bottom--md', styles.supportedPackagesHeader)}
          />
          <ul className={clsx('clean-list', styles.packageList)}>
            {filteredUsers.map((user) => (
              <PackageCard key={user.name} user={user}/>
            ))}
          </ul>
        </div>
      )}
    </section>
  );
}

const favoriteUsers = sortedUsers.filter((user) =>
  !user.tags.includes('planned'),
);
const otherUsers = sortedUsers.filter(
  (user) => user.tags.includes('planned'),
);


export default function PackagePage(): JSX.Element {
  return (
    <Layout title={TITLE} description={DESCRIPTION}>
      <main className="margin-vert--lg">
        <PackagesHeader/>
        <PackagesFilters/>
        <div
          style={{display: 'flex', marginLeft: 'auto', justifyContent: 'center'}}
          className="container">
          <SearchBar/>
        </div>
        <Packages/>
      </main>
    </Layout>
  );
}
