chrome.runtime.onInstalled.addListener(() => {
  console.log("PassVault Auto-Fill extension installed.");
});

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "getAuthToken") {
    chrome.storage.local.get(["authToken"], function (result) {
      sendResponse({ token: result.authToken });
    });
    return true; // Will respond asynchronously
  }
  if (request.action === "storeAuthToken") {
    chrome.storage.local.set({ authToken: request.token }, () => {});
    return true;
  }
  if (request.action === "removeAuthToken") {
    chrome.storage.local.remove("authToken", () => {});
    return true;
  }
});

chrome.action.onClicked.addListener((tab) => {
  chrome.scripting.executeScript({
    target: { tabId: tab.id },
    files: ["content.js"],
  });
});
