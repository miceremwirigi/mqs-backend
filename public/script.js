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