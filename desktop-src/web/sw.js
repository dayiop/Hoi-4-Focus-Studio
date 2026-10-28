// Works for a user/repository.github.io/project/ prefix without hardcoded paths.
const CACHE = 'hoi4-focus-studio-v4-20261028';
const ASSETS = ['./', './index.html', './manifest.webmanifest', './icon.svg'];
self.addEventListener('install', event => {
  event.waitUntil(caches.open(CACHE).then(cache => cache.addAll(ASSETS)).then(() => self.skipWaiting()));
});
self.addEventListener('activate', event => {
  event.waitUntil(caches.keys().then(keys => Promise.all(keys.filter(k => k.startsWith('hoi4-focus-studio-') && k !== CACHE).map(k => caches.delete(k)))).then(() => self.clients.claim()));
});
self.addEventListener('fetch', event => {
  if(event.request.method !== 'GET' || new URL(event.request.url).origin !== self.location.origin) return;
  event.respondWith(fetch(event.request).then(response => {
    if(response.ok){const copy=response.clone();caches.open(CACHE).then(c => c.put(event.request,copy)).catch(()=>{});}
    return response;
  }).catch(() => caches.match(event.request).then(r => r || (event.request.mode==='navigate' ? caches.match('./index.html') : Response.error()))));
});
