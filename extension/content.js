chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "fillPassword") {
    const { username, password } = request.data;

    // Here you would identify the username and password fields on the page
    // and fill them. This example assumes input fields with id 'username' and 'password'.
    const usernameField =
      document.querySelector('input[name="username"]') ||
      document.querySelector('input[type="email"]');
    const passwordField =
      document.querySelector('input[name="password"]') ||
      document.querySelector('input[type="password"]');

    if (usernameField && passwordField) {
      usernameField.value = username;
      passwordField.value = password;
    }

    sendResponse({ success: true });
  }
});

window.addEventListener(
  "message",
  function (event) {
    if (event.source !== window) return;
    if (event.source.origin !== "https://www.passvault.fun") {
      return;
    }
    if (event.data.type && event.data.type === "PASSVAULT_TOKEN") {
      const token = event.data.token;
      chrome.runtime.sendMessage(
        { action: "storeAuthToken", token: token },
        () => {},
      );
    }

    if (event.data.type && event.data.type === "REMOVE_PASSVAULT_TOKEN") {
      chrome.runtime.sendMessage({ action: "removeAuthToken" }, () => {});
    }
  },
  false,
);
