// See: https://serviceworke.rs/

self.addEventListener('push', function(event) {
    console.log('[Service Worker] Push Received.');
    const payload = event.data ? event.data.json() : 'no payload';
    self.registration.showNotification(payload.title, {
        body: payload.body,
        icon: payload.icon,
    })
});