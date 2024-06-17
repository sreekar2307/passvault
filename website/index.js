// global variables
var passwords = [];
var hiddenElements;
var currentPage = 1;
var itemsPerPage = 10;
var fetchedPages = new Set();
var searchQuery = "";
var currentPasswordRequest;
var beEndpoint = "https://www.passvault.fun/api/v1";
var beEndpointV2 = "https://www.passvault.fun/api/v2";
var passKeyCreationAllowed = false;

function loginNavClick() {
  $("#login").removeClass("hidden");
  $("#create-user").addClass("hidden");
}

function showLoader() {
  $("#loader").removeClass("d-none");
}

function hideLoader() {
  $("#loader").addClass("d-none");
}

function closeUpdateDrawer() {
  $("#updateDrawer").removeClass("show");
}

function closeStoreDrawer() {
  $("#storeDrawer").removeClass("show");
}

function setPage(page) {
  currentPage = page;
  if (!fetchedPages.has(page)) {
    fetchPasswords(page, function () {
      displayCurrentPage();
    });
  } else {
    displayCurrentPage();
  }
}

function updatePaginationControls() {
  var paginationControls = $("#paginationControls");
  paginationControls.empty();

  for (var i = 1; i <= passwords.length + 1; i++) {
    var activeClass = i === currentPage ? "active" : "";
    var pageItem = `
      <li class="page-item ${activeClass}">
          <a class="page-link" href="#" onclick="setPage(${i})">${i}</a>
      </li>`;
    paginationControls.append(pageItem);
  }
}

function signUpNavClick() {
  $("#create-user").removeClass("hidden");
  $("#login").addClass("hidden");
}

function displayPasswords(passwords) {
  var tableBody = $("#passwordTableBody");
  tableBody.empty();
  passwords.forEach(function (password) {
    var row = `<tr>
      <td>${password.website}</td>
      <td>${password.username.String}</td>
      <td>${password.email.String}</td>
      <td>
          <div class="d-flex">
              <button class="btn btn-sm mr-2" onclick="showUpdateDrawer(${password.ID})">
                  <i class="fas fa-edit"></i>
              </button>
              <button class="btn btn-sm" onclick="submitDeletePassword(${password.ID})">
                  <i class="fas fa-trash-alt"></i>
              </button>
          </div>
      </td>
      </tr>`;
    tableBody.append(row);
  });
}

function displayCurrentPage() {
  if (passwords) {
    var passwordsToDisplay = passwords[currentPage - 1];
    if (passwordsToDisplay) {
      displayPasswords(passwordsToDisplay);
      updatePaginationControls();
    }
  }
}

function fetchPasswords(page, callback) {
  // Abort the previous request if it exists
  if (currentPasswordRequest) {
    currentPasswordRequest.abort();
  }

  var token = localStorage.getItem("authToken");
  var offset = (page - 1) * itemsPerPage;
  currentPasswordRequest = $.ajax({
    url: beEndpoint + "/passwords",
    method: "GET",
    headers: {
      Authorization: "Bearer " + token,
    },
    data: { offset: offset, limit: itemsPerPage, query: searchQuery }, // Assuming your backend supports pagination parameters
    success: function (response) {
      if (response.data && response.data.length > 0) {
        passwords.push(response.data);
        fetchedPages.add(page);
      } else {
        fetchedPages.add(page);
      }
      callback();
    },
    error: function (jqXHR, textStatus) {
      handleBEAPIError(jqXHR, textStatus);
    },
    complete: function () {
      currentPasswordRequest = null; // Clear the current request
    },
  });
}

function showUpdateDrawer(id) {
  var password = passwords[currentPage - 1].find((p) => p.ID === id);
  if (password) {
    $("#updateID").val(password.ID), $("#updateName").val(password.name.String);
    $("#updateWebsite").val(password.website);
    $("#updateUsername").val(password.username.String);
    $("#updatePassword").val(password.password);
    $("#updateEmail").val(password.email.String);
    $("#updateNotes").val(password.notes.String);
    $("#updateDrawer").addClass("show");
  }
}

function showStoreDrawer() {
  $("#storeDrawer").addClass("show");
}

function showImportDrawer() {
  $("#importDrawer").addClass("show");
}

function checkGenerateConditions() {
  var symbols = $("#symbols").is(":checked");
  var alphabets = $("#alphabets").is(":checked");
  var numbers = $("#numbers").is(":checked");

  if (symbols || alphabets || numbers) {
    $("#generateButton").prop("disabled", false);
  } else {
    $("#generateButton").prop("disabled", true);
  }
}

function submitPasswordForm() {
  var passwordData = JSON.stringify({
    name: $("#storeName").val(),
    website: $("#storeWebsite").val(),
    username: $("#storeUsername").val(),
    password: $("#storePassword").val(),
    email: $("#storeEmail").val(),
    notes: $("#storeNotes").val(),
  });

  // Call the storePassword function with the password data
  storePassword(passwordData);
}

function resetTable() {
  fetchedPages.clear();
  passwords = [];
  setPage(1);
}

function storePassword(passwordData) {
  var token = localStorage.getItem("authToken");

  $.ajax({
    url: beEndpoint + "/passwords", // Replace with your backend API URL
    method: "POST",
    headers: {
      Authorization: "Bearer " + token,
    },
    contentType: "application/json",
    data: passwordData,
    success: function (response) {
      closeStoreDrawer();
      resetTable();
    },
    error: function (jqXHR, textStatus) {
      // Handle error response if needed
      handleBEAPIError(jqXHR, textStatus);
    },
  });
}

$("#storePasswordForm").submit(function (e) {
  e.preventDefault();
  // Placeholder for store functionality
  submitPasswordForm();
});

$("#createUserForm").submit(function (e) {
  e.preventDefault();
  submitCreateuserForm();
});

function submitCreateuserForm() {
  var createUserData = {
    email: $("#createUserEmail").val(),
    name: $("#createUserUsername").val(),
  };
  createUser(createUserData);
}

function createUser(createUserData) {
  // getCaptchaToken(function (token) {
  //   createUserData.token = token;
  //
  // });
  $.ajax({
    url: beEndpointV2 + "/begin/register",
    method: "POST",
    contentType: "application/json",
    data: JSON.stringify(createUserData),
    success: async function (response) {
      if (response.token) {
        postLogin(response.token);
      }
      let dataDecoded = response.credOptions;
      dataDecoded.publicKey.challenge = base64ToArrayBuffer(
        dataDecoded.publicKey.challenge,
      );
      dataDecoded.publicKey.user.id = base64ToArrayBuffer(
        dataDecoded.publicKey.user.id,
      );
      const credential = await navigator.credentials.create(dataDecoded);
      console.log("printing credentials", credential);
      sendCredentialToBE(credential, response.sessionID);
    },
    error: function (jqXHR, textStatus) {
      // Handle error response if needed
      handleBEAPIError(jqXHR, textStatus);
    },
  });
}

// Function to convert ArrayBuffer to Base64
function arrayBufferToBase64(buffer) {
  // Create a Uint8Array from the ArrayBuffer
  let binary = "";
  const bytes = new Uint8Array(buffer);
  const len = bytes.byteLength;

  // Convert each byte to a binary string
  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }

  // Convert the binary string to a Base64 encoded string
  const base64String = window.btoa(binary);

  // Make the Base64 string URL-safe by replacing non-URL-safe characters
  const urlSafeBase64 = base64String
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/, "");

  return urlSafeBase64;
}

function sendCredentialToBE(credential, sessionID) {
  $.ajax({
    url: beEndpointV2 + "/finish/register" + "?session_id=" + sessionID,
    method: "POST",
    contentType: "application/json",
    data: JSON.stringify({
      rawId: arrayBufferToBase64(credential.rawId),
      id: credential.id,
      type: credential.type,
      response: {
        clientDataJSON: arrayBufferToBase64(credential.response.clientDataJSON),
        attestationObject: arrayBufferToBase64(
          credential.response.attestationObject,
        ),
      },
      authenticatorAttachment: credential.authenticatorAttachment,
    }),
    success: async function (response) {
      if (response.token) {
        postLogin(response.token);
      }
    },
    error: function (jqXHR, textStatus) {
      // Handle error response if needed
      handleBEAPIError(jqXHR, textStatus);
    },
  });
}

function base64Decode(input) {
  // Replace non-url compatible chars with base64 standard chars
  input = input.replace(/-/g, "+").replace(/_/g, "/");

  // Pad out with standard base64 required padding characters
  var pad = input.length % 4;
  if (pad) {
    if (pad === 1) {
      throw new Error(
        "InvalidLengthError: Input base64url string is the wrong length to determine padding",
      );
    }
    input += new Array(5 - pad).join("=");
  }

  return atob(input);
}

function base64ToArrayBuffer(base64) {
  return stringAsBytes(base64Decode(base64));
}

function stringAsBytes(str) {
  var bytes = new Uint8Array(str.length);
  for (var i = 0; i < str.length; i++) {
    bytes[i] = str.charCodeAt(i);
  }
  return bytes.buffer;
}

$("#importPasswordForm").submit(function (e) {
  alert("Passwords imported");
});

function submitDeletePassword(id) {
  deletePassword(id);
}

function deletePassword(id) {
  var token = localStorage.getItem("authToken");

  $.ajax({
    url: beEndpoint + "/passwords/" + id,
    method: "DELETE",
    contentType: "application/json",
    headers: {
      Authorization: "Bearer " + token,
    },
    success: function (response) {
      resetTable();
    },
    error: function (jqXHR, textStatus) {
      // Handle error response if needed
      handleBEAPIError(jqXHR, textStatus);
    },
  });
}

$("#generatePasswordForm").submit(function (e) {
  e.preventDefault();
  var size = $("#passwordSize").val();
  var symbols = $("#symbols").is(":checked");
  var alphabets = $("#alphabets").is(":checked");
  var numbers = $("#numbers").is(":checked");
  generatePassword(size, symbols, alphabets, numbers);
});

function generatePassword(size, symbols, alphabets, numbers) {
  var token = localStorage.getItem("authToken");

  $.ajax({
    url: beEndpoint + "/generate/passwords",
    method: "POST",
    contentType: "application/json",
    headers: {
      Authorization: "Bearer " + token,
    },
    data: JSON.stringify({
      size: parseInt(size),
      excludeSymbols: !symbols,
      excludeDigits: !numbers,
      excludeAlphabets: !alphabets,
    }),
    success: function (response) {
      $("#generatedPassword").text(response.password);
    },
    error: function (jqXHR, textStatus) {
      // Handle error response if needed
      handleBEAPIError(jqXHR, textStatus);
    },
  });
}

function submitUpdatePasswordForm() {
  var passwordData = JSON.stringify({
    name: $("#updateName").val(),
    website: $("#updateWebsite").val(),
    username: $("#updateUsername").val(),
    password: $("#updatePassword").val(),
    email: $("#updateEmail").val(),
    notes: $("#updateNotes").val(),
  });

  // Call the updatePassword function with the password data
  updatePassword(passwordData, $("#updateID").val());
}

function updatePassword(passwordData, id) {
  var token = localStorage.getItem("authToken");

  $.ajax({
    url: beEndpoint + "/passwords/" + id, // Replace with your backend API URL
    method: "PUT",
    contentType: "application/json",
    headers: {
      Authorization: "Bearer " + token,
    },
    data: passwordData,
    success: function (response) {
      closeUpdateDrawer();
      resetTable();
    },
    error: function (jqXHR, textStatus) {
      // Handle error response if needed
      handleBEAPIError(jqXHR, textStatus);
    },
  });
}

$("#updatePasswordForm").submit(function (e) {
  e.preventDefault();
  submitUpdatePasswordForm();
});

function showGenerateOptions() {
  $("#generateOptions").toggleClass("hidden");
}

function logout() {
  localStorage.removeItem("authToken");
  if (hiddenElements) {
    hiddenElements.addClass("hidden");
  }
  sendRemoveTokenToExtension();
  $("#create-user").addClass("hidden");
  $("#login").removeClass("hidden");
  $("#loginNavItem").removeClass("hidden");
  $("#signUpNavItem").removeClass("hidden");
  $("#generateOptions").removeClass("hidden");
}

function handleBEAPIError(jqXHR, textStatus) {
  if (jqXHR.status == 401) {
    logout();
  } else if (textStatus !== "abort") {
    alert(jqXHR.status);
  }
}

function postLogin(token) {
  localStorage.setItem("authToken", token);
  sendTokenToExtension(token);
  hiddenElements = $(".hidden");
  hiddenElements.removeClass("hidden");
  $("#create-user, #login").addClass("hidden");
  $("#loginNavItem").addClass("hidden");
  $("#signUpNavItem").addClass("hidden");
  $("#generateOptions").addClass("hidden");
  setPage(1);
}

function sendTokenToExtension(token) {
  window.postMessage({ type: "PASSVAULT_TOKEN", token: token }, "*");
}

function sendRemoveTokenToExtension() {
  window.postMessage({ type: "REMOVE_PASSVAULT_TOKEN" }, "*");
}

function searchPasswords(query) {
  passwords = [];
  fetchedPages.clear();
  currentPage = 1;
  setPage(1);
}

// function getCaptchaToken(callback) {
//   grecaptcha.ready(function () {
//     grecaptcha
//       .execute(googleCaptchaToken, {
//         action: "submit",
//       })
//       .then(function (token) {
//         callback(token);
//       })
//       .catch(function (error) {
//         alert(error);
//       });
//   });
// }

$(document).ready(function () {
  if (localStorage.getItem("authToken")) {
    postLogin(localStorage.getItem("authToken"));
  }

  if (
    window.PublicKeyCredential &&
    PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable &&
    PublicKeyCredential.isConditionalMediationAvailable
  ) {
    // Check if user verifying platform authenticator is available.
    Promise.all([
      PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable(),
      PublicKeyCredential.isConditionalMediationAvailable(),
    ]).then((results) => {
      if (results.every((r) => r === true)) {
        passKeyCreationAllowed = true;
      }
    });
  }

  const toggleLoginIconSpan = $("#toggleLoginPassword");
  const toggleCreateUserIconSpan = $("#toggleCreateUserPassword");
  const toggleCreateUserConfirmIconSpan = $("#toggleCreateUserConfirmPassword");
  const toggleUpdateIconSpan = $("#toggleUpdatePassword");
  const toggleStoreIconSpan = $("#toggleStorePassword");

  toggleLoginIconSpan.on("click", function () {
    const passwordField = $("#loginPassword");
    const toggleIcon = $("#toggleIconLoginPassword");
    togglePassWordEye(passwordField, toggleIcon);
  });

  toggleCreateUserIconSpan.on("click", function () {
    const passwordField = $("#createUserPassword");
    const toggleIcon = $("#toggleIconCreateUserPassword");
    togglePassWordEye(passwordField, toggleIcon);
  });

  toggleCreateUserConfirmIconSpan.on("click", function () {
    const passwordField = $("#createUserConfirmPassword");
    const toggleIcon = $("#toggleIconCreateUserConfirmPassword");
    togglePassWordEye(passwordField, toggleIcon);
  });

  toggleUpdateIconSpan.on("click", function () {
    const passwordField = $("#updatePassword");
    const toggleIcon = $("#toggleIconUpdatePassword");
    togglePassWordEye(passwordField, toggleIcon);
  });

  toggleStoreIconSpan.on("click", function () {
    const passwordField = $("#storePassword");
    const toggleIcon = $("#toggleIconStorePassword");
    togglePassWordEye(passwordField, toggleIcon);
  });

  function togglePassWordEye(passwordField, toggleIconField) {
    if (passwordField.attr("type") === "password") {
      passwordField.attr("type", "text");
      toggleIconField.removeClass("fa-eye").addClass("fa-eye-slash");
    } else {
      passwordField.attr("type", "password");
      toggleIconField.removeClass("fa-eye-slash").addClass("fa-eye");
    }
  }

  $("#importPasswordForm label").click(function (e) {
    e.preventDefault();
    $("#importCSV").click();
  });

  $("#importCSV").change(function () {
    var formData = new FormData();
    var file = $("#importCSV")[0].files[0];

    if (file) {
      formData.append("passwords", file);
      showLoader();
      var token = localStorage.getItem("authToken");
      $.ajax({
        url: beEndpoint + "/import/passwords", // Replace with your backend API URL
        method: "POST",
        headers: {
          Authorization: "Bearer " + token,
        },
        data: formData,
        processData: false,
        contentType: false,
        success: function (response) {
          resetTable();
        },
        error: function (jqXHR, textStatus) {
          handleBEAPIError(jqXHR, textStatus);
        },
        complete: function () {
          hideLoader();
        },
      });
    }
  });

  $("#searchInput").on("input", function () {
    searchQuery = $(this).val().toLowerCase();
    searchPasswords(searchQuery);
  });

  $("#loginForm").submit(function (event) {
    event.preventDefault();
    var email = $("#loginEmail").val();

    $.ajax({
      url: beEndpointV2 + "/begin/login",
      method: "POST",
      contentType: "application/json",
      data: JSON.stringify({
        email: email,
      }),
      success: async function (response) {
        if (response.credAssertion) {
          let dataDecoded = response.credAssertion;
          dataDecoded.publicKey.challenge = base64ToArrayBuffer(
            dataDecoded.publicKey.challenge,
          );
          for (
            let i = 0;
            i < dataDecoded.publicKey.allowCredentials.length;
            i++
          ) {
            dataDecoded.publicKey.allowCredentials[i].id = base64ToArrayBuffer(
              dataDecoded.publicKey.allowCredentials[i].id,
            );
          }

          await finishLogin(
            response.credAssertion.publicKey,
            response.sessionID,
          );
        }
      },
      error: function () {
        alert("Login failed. Please try again.");
      },
    });
  });
});

async function finishLogin(publicKeyOptions, sessionID) {
  const credential = await navigator.credentials.get({
    publicKey: publicKeyOptions,
  });

  console.log("printing credentials", credential);
  let credentialParsed = {
    id: credential.id,
    type: credential.type,
    authenticatorAttachment: credential.authenticatorAttachment,
    rawId: arrayBufferToBase64(credential.rawId),
    response: {
      clientDataJSON: arrayBufferToBase64(credential.response.clientDataJSON),
      authenticatorData: arrayBufferToBase64(
        credential.response.authenticatorData,
      ),
      userHandle: arrayBufferToBase64(credential.response.userHandle),
      signature: arrayBufferToBase64(credential.response.signature),
    },
  };

  $.ajax({
    url: beEndpointV2 + "/finish/login" + "?session_id=" + sessionID,
    method: "POST",
    contentType: "application/json",
    data: JSON.stringify(credentialParsed),
    success: function (response) {
      if (response.token) {
        postLogin(response.token);
      } else {
        alert("Login failed. Please try again.");
      }
    },
    error: function () {
      alert("Login failed. Please try again.");
    },
  });
}
