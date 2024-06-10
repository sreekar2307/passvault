chrome.runtime.onInstalled.addListener(() => {
    console.log('PassVault Auto-Fill extension installed.');
});

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === 'getAuthToken') {
        console.log("got message for getAuthToken", request)
        chrome.storage.local.get(['authToken'], function(result) {
            console.log("sending the token ", result)
            sendResponse({ token: result.authToken });
        });
        return true; // Will respond asynchronously
    }
});

chrome.action.onClicked.addListener((tab) => {
    chrome.scripting.executeScript({
        target: { tabId: tab.id },
        files: ['content.js']
    });
});

