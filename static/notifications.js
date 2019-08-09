// See: https://serviceworke.rs/

function urlB64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
        .replace(/\-/g, '+')
        .replace(/_/g, '/');

    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);

    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
}

navigator.serviceWorker.register('sw.js');

navigator.serviceWorker.ready
    .then(function(registration) {

        return registration.pushManager.getSubscription()
            .then(async function(subscription) {

                if (subscription) {
                    return subscription;
                }

                const response = await fetch('./publicKey');
                const vapidPublicKey = await response.text();

                const convertedVapidKey = urlB64ToUint8Array(vapidPublicKey);


                return registration.pushManager.subscribe({
                    userVisibleOnly: true,
                    applicationServerKey: convertedVapidKey
                });
            });
    }).then(function(subscription) {
        fetch('./subscribe', {
            method: 'post',
            headers: {
                'Content-type': 'application/json'
            },
            body: JSON.stringify(subscription),
        });

    });