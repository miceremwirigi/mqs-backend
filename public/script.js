async function fetchWithAuth(url, options = {}) {
  const defaultHeaders = {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer ' + localStorage.getItem('jwt_token'),
  };
  options.headers = { ...defaultHeaders, ...options.headers };

  const res = await fetch(url, options);
  if (!res.ok) {
    redirectToLoginIfUnauthorized(res);
    // Try to parse error JSON, fallback to status text
    let errorMsg = `HTTP Error: Status ${res.status}`;
    try {
      const errJson = await res.json();
      errorMsg = errJson.message || errorMsg;
    } catch (e) {
    }
    throw new Error(errorMsg);
  }
  return await res.json();
}

function redirectToLoginIfUnauthorized(response) {
  if (response.status === 401) {
      localStorage.removeItem('jwt_token');
      window.location.href = '/login.html';
  }
}

var allHospitals = new Array;
var allEngineers = new Array;
var allEquipment = new Array;
var allServices = new Array;

// set form visibility toggle
function toggleFormVisibility(formId) {
    const form = document.getElementById(formId);
    form.style.display = form.style.display === 'none' || form.style.display === '' ? 'flex' : 'none';
    form.scrollIntoView({ behavior: "smooth", block: "center" });
}

// Initially hide the input forms
function hideInputForms(){
    var addForms = document.getElementsByClassName("add-content-form");
    for (var i=0; i<addForms.length; i++){
        addForms[i].style.display = 'none';
    }
}
setTimeout(hideInputForms,1)

function showSuccessMessage() {
  const popup = document.getElementById('success-popup');
  popup.style.display = 'block';
  setTimeout(function() {
    popup.style.display = 'none';
  }, 3000); // Message disappears after 3 seconds (3000 milliseconds)
}

function showFailedMessage() {
  const popup = document.getElementById('failed-popup');
  popup.style.display = 'block';
  setTimeout(function() {
    popup.style.display = 'none';
  }, 3000); // Message disappears after 3 seconds (3000 milliseconds)
}

// Logout utility (call this on logout)
function handleLogout() {
    localStorage.removeItem('jwt_token');
    window.location.href = "login.html";
}

// Populate department select
async function populateDepartmentSelect() {
    const sel = document.getElementById('equipment-department-id');
    sel.innerHTML = '';
    const res = await fetchWithAuth('/api/departments');
    (res.data || []).forEach(dept => {
        const opt = document.createElement('option');
        opt.value = dept.ID || dept.id;
        opt.textContent = dept.Name || dept.name;
        sel.appendChild(opt);
    });
}

document.addEventListener("DOMContentLoaded", function () {
  const dropdown = document.querySelector('.menu-bar .dropdown');
  const dropbtn = dropdown.querySelector('.dropbtn');
  const dropdownContent = dropdown.querySelector('.dropdown-content');
  const closeBtn = document.getElementById('dropdown-close-btn');
  dropbtn.addEventListener('click', function (e) {
    e.preventDefault();
    dropdown.classList.toggle('active');
  });
  // Hide dropdown-content when a menu link is clicked (on all screens)
  const menuLinks = dropdownContent.querySelectorAll('.menu-link');
  menuLinks.forEach(link => {
    link.addEventListener('click', function () {
      dropdown.classList.remove('active');
    });
  });
  // Hide dropdown-content when X is clicked
  if (closeBtn) {
    closeBtn.addEventListener('click', function () {
      dropdown.classList.remove('active');
    });
  }
  // Hide dropdown-content when clicking outside
  document.addEventListener('click', function (e) {
    if (!dropdown.contains(e.target)) {
      dropdown.classList.remove('active');
    }
  });
});



