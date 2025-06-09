var token = localStorage.getItem('jwt_token');

function fetchWithAuth(url, options = {}) {
  const defaultHeaders = {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('jwt_token'),
  };
  options.headers = { ...defaultHeaders, ...options.headers };

  return fetch(url, options)
      .then(res => {
          if (!res.ok) {
              redirectToLoginIfUnauthorized(res);
              throw new Error(`HTTP Error: Status ${res.status}`);
          }
          return res.json();
      });
}

function redirectToLoginIfUnauthorized(response) {
  if (response.status === 401) {
      localStorage.removeItem('jwt_token');
      window.location.href = '/login.html';
  }
}   
document.addEventListener("DOMContentLoaded", function() {
        const menuLinks = document.querySelectorAll('.menu-link');
        const sections = document.querySelectorAll('.main-section');
        function showSection(id) {
          sections.forEach(sec => sec.style.display = 'none');
          const target = document.getElementById(id);
          if (target) {
            target.style.display = 'block';
            localStorage.setItem('activeSection', id);
          }
        }
        menuLinks.forEach(link => {
          link.addEventListener('click', function(e) {
            e.preventDefault();
            showSection(this.dataset.section);
          });
        });
        // Show the first section by default
        const activeSection = localStorage.getItem('activeSection');
        if (activeSection && document.getElementById(activeSection)) {
          sections.forEach(sec => sec.style.display = 'none');
          document.getElementById(activeSection).style.display = 'block';
        } else {
          if (sections[0]) {
            sections[0].style.display = 'block';
          }
        }
      });

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

// Show/hide add department form
const showAddDeptBtn = document.getElementById('show-add-department-form-btn');
const addDeptForm = document.getElementById('add-department-form-modal');
if (showAddDeptBtn && addDeptForm) {
    showAddDeptBtn.onclick = () => {
        addDeptForm.style.display = addDeptForm.style.display === 'none' ? '' : 'none';
    };
}

const submitDeptBtn = document.getElementById('submit-department-modal');
if (submitDeptBtn && addDeptForm) {
    submitDeptBtn.onclick = async () => {
        const nameInput = document.getElementById('modal-department-name');
        if (!nameInput) return;
        const name = nameInput.value;
        if (!name) return alert('Enter department name');
        await fetchWithAuth('/api/departments', {
            method: 'POST',
            body: JSON.stringify({ name }),
            headers: { 'Content-Type': 'application/json' }
        });
        await populateDepartmentSelect();
        addDeptForm.style.display = 'none';
        nameInput.value = '';
    };
}
