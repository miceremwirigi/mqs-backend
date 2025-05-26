// set form visibility toggle
function toggleFormVisibility(formId) {
    const form = document.getElementById(formId);
    form.style.display = form.style.display === 'none' || form.style.display === '' ? 'flex' : 'none';
}

// Initially hide the input forms
function hideInputForms(){
    var addForms = document.getElementsByClassName("add-content-form");
    for (var i=0; i<addForms.length; i++){
        addForms[i].style.display = 'none';
    }
}
setTimeout(hideInputForms,1)
