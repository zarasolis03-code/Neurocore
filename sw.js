self.addEventListener('install', (e) => {
  console.log('Neurocore Service Worker Installed');
});
self.addEventListener('fetch', (e) => {
  e.respondWith(fetch(e.request));
});
