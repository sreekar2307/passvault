// global variables
var passwords = [];
var userAuthToken;
var isAuthenticated = false;
var hiddenElements;
var currentPage = 1;
var itemsPerPage = 10;
var fetchedPages = new Set();
var searchQuery = "";
var currentPasswordRequest;
var beEndpoint = "https://www.passvault.fun/api/v1"

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
      <td>${password.password}</td>
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
    $("#updateID").val(password.ID),
      $("#updateName").val(password.name.String);
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

$('#createUserForm').submit(function (e) {
  e.preventDefault();
  submitCreateuserForm();
});

function submitCreateuserForm() {
  var createUserData = JSON.stringify({
    email: $("#createUserEmail").val(),
    name: $("#createUserUsername").val(),
    password: $("#createUserPassword").val(),
    confirmPassword: $("#createUserConfirmPassword").val(),
  })
  createUser(createUserData)
}

function createUser(createUserData) {
  $.ajax({
    url: beEndpoint + "/users",
    method: "POST",
    contentType: "application/json",
    data: createUserData,
    success: function (response) {
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
  localStorage.setItem("authToken", token); // Store the token in localStorage
  hiddenElements = $(".hidden");
  hiddenElements.removeClass("hidden");
  $("#create-user, #login").addClass("hidden");
  $("#loginNavItem").addClass("hidden");
  $("#signUpNavItem").addClass("hidden");
  $("#generateOptions").addClass("hidden");
  setPage(1);
}

function searchPasswords(query) {
  passwords = [];
  fetchedPages.clear();
  currentPage = 1;
  setPage(1);
}

$(document).ready(function () {
  if (localStorage.getItem("authToken")) {
    postLogin(localStorage.getItem("authToken"));
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
    var password = $("#loginPassword").val();

    $.ajax({
      url: beEndpoint + "/login/users", // Replace with your backend API URL
      method: "POST",
      contentType: "application/json",
      data: JSON.stringify({
        email: email,
        password: password,
      }),
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
  });
});