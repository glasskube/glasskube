(() => {
  const getColorSchemeQuery = () =>
    window.matchMedia('(prefers-color-scheme: dark)');
  const getPreferredTheme = () =>
    getColorSchemeQuery().matches ? 'dark' : 'light';
  const setPreferredTheme = () =>
    document.body.setAttribute('data-bs-theme', getPreferredTheme());
  setPreferredTheme();
  getColorSchemeQuery().addEventListener('change', () => setPreferredTheme());
})();

(() => {
  const dismissed = sessionStorage.getItem('cloud-info-dismissed') === 'true';
  if (!dismissed) {
    document.body
      .querySelector('#cloud-info')
      .classList.remove('visually-hidden');
  }
  document.body
    .querySelector('#cloud-info-close')
    .addEventListener('click', () => {
      sessionStorage.setItem('cloud-info-dismissed', true);
    });
})();

window.advancedOptions = function (currentContext) {
  return localStorage.getItem('advancedOptions_' + currentContext) === 'true';
};

window.sseConnected = false;
function showDisconnectedToast(show) {
  const elem = document.getElementById('disconnected-toast');
  if (!elem) {
    return;
  }
  if (!show && elem.classList.contains('show')) {
    document.getElementById('disconnected-toast').classList.remove('show');
  } else if (show && !elem.classList.contains('show')) {
    document.getElementById('disconnected-toast').classList.add('show');
  }
}
document.addEventListener('htmx:sseError', function (evt) {
  console.log('htmx:sseError', evt);
  window.sseConnected = false;
  showDisconnectedToast(true);
});
document.addEventListener('htmx:sseClose', function (evt) {
  console.log('htmx:sseClose', evt);
  window.sseConnected = false;
  setTimeout(() => {
    if (!window.sseConnected) {
      showDisconnectedToast(true);
    }
  }, 1000);
});
document.addEventListener('htmx:sseOpen', function (evt) {
  console.log('htmx:sseOpen', evt);
  window.sseConnected = true;
  showDisconnectedToast(false);
});

window.giscusReported = false;
function handleGiscusMessage(ev) {
  if (window.giscusReported) return;
  if (ev.origin !== 'https://giscus.app') return;
  if (!(typeof ev.data === 'object' && ev.data.giscus)) return;

  const giscusData = ev.data.giscus;
  if (giscusData['discussion'] && giscusData['viewer']) {
    const username = giscusData['viewer']['login'];
    if (username.includes('giscus') && username.includes('bot')) {
      return;
    }
    const githubUrl = giscusData['viewer']['url'];
    const formData = new FormData();
    formData.append('githubUrl', githubUrl);
    fetch('', {
      method: 'POST',
      body: formData,
    });
    window.giscusReported = true;
  }
}
window.addEventListener('message', handleGiscusMessage);
