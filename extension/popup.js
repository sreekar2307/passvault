function copyToClipboard(text) {
  const tempInput = $("<input>");
  $("body").append(tempInput);
  tempInput.val(text).select();
  document.execCommand("copy");
  tempInput.remove();
}

const limit = 10;
let offset = 0;
let query = "";
let passwords = [];
let reachedEnd = false;

function renderPasswords(passwords) {
  const $passwordList = $("#password-list");
  $passwordList.empty(); // Clear existing passwords

  passwords.forEach((password) => {
    const $passwordItem = $(`
      <li class="list-group-item">
      ${password.username.String}
      ${password.website}
      <div>
          <button class="btn btn-link btn-sm copy-username" data-username="${password.username.String}">
              <i class="fa-solid fa-user"></i>
          </button>
          <button class="btn btn-link btn-sm copy-email" data-email="${password.email.String}">
              <i class="fas fa-envelope"></i>
          </button>
          <button class="btn btn-link btn-sm copy-password" data-password="${password.password}">
              <i class="fas fa-key"></i>
          </button>
      </div>
  </li>
            `);
    $passwordList.append($passwordItem);
  });

  $(".copy-username").on("click", function () {
    const username = $(this).data("username");
    copyToClipboard(username);
  });

  $(".copy-email").on("click", function () {
    const email = $(this).data("email");
    copyToClipboard(email);
  });

  $(".copy-password").on("click", function () {
    const password = $(this).data("password");
    copyToClipboard(password);
  });
}

$("#search").on("input", async function () {
  query = $(this).val().toLowerCase();
  offset = 0;
  reachedEnd = false;
  passwords = await fetchPasswords(limit, offset, query);
  renderPasswords(passwords);
});

async function getAuthToken() {
  return new Promise((resolve) => {
    chrome.runtime.sendMessage({ action: "getAuthToken" }, function (response) {
      resolve(response.token);
    });
  });
}

async function fetchPasswords(limit, offset, query) {
  await chrome.tabs.query({ active: true, currentWindow: true });
  const token = await getAuthToken();
  const url = new URL("https://www.passvault.fun/api/v1/passwords");
  url.searchParams.append("limit", limit.toString());
  url.searchParams.append("offset", offset.toString());
  url.searchParams.append("query", query);
  const response = await fetch(url, {
    method: "GET",
    headers: {
      Authorization: "Bearer " + token,
    },
  });
  const data = await response.json();
  return data.data;
}

async function handleScroll() {
  const $container = $("#password-container");
  const scrollTop = $container.scrollTop();
  const containerHeight = $container.height();
  const contentHeight = $container[0].scrollHeight;

  if (scrollTop + containerHeight >= contentHeight - 100) {
    offset += limit;
    if (!reachedEnd) {
      currentPasswords = await fetchPasswords(limit, offset, query);
      if (currentPasswords.length === 0) {
        reachedEnd = true;
      } else {
        passwords = passwords.concat(currentPasswords);
        renderPasswords(passwords);
      }
    }
  }
}

$(document).ready(function () {
  $("#password-container").on("scroll", handleScroll);
});

// Trigger the fetch passwords on load if token is available
chrome.runtime.sendMessage(
  { action: "getAuthToken" },
  async function (response) {
    const token = response.token;
    if (token) {
      passwords = await fetchPasswords(limit, offset, query);
      renderPasswords(passwords);
      $("#searchFormGroup").removeClass("d-none");
      $("#password-container").removeClass("d-none");
    } else {
      $("#searchFormGroup").addClass("d-none");
      $("#password-container").addClass("d-none");
    }
  },
);
