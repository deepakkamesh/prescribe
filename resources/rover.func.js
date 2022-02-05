

$(document).ready(function() {
  var errorContainer = document.querySelector('#error-popup');

  // Do a system health check on system metrics on a internal.
  var status_check = function() {
    $.post('/api/status', "", function(data, status) {
      if (data.Err != '') {
        console.log(data.Err);
        errorContainer.MaterialSnackbar.showSnackbar({message : data.Err});
        return;
      }
    });
  };
  setInterval(status_check, 5000);

  var currPrescription = ""; // Filename of current prescription.

  // Generate PDF button handler.
  document.querySelector('#generate_pdf').addEventListener('click', function() {
    name = $('#name').val();
    age_sex = $('#age_sex').val();
    prescription = $('#prescription').val();

    // Do some basic validation.
    if ($.trim(name) == "" || $.trim(age_sex) == "" ||
        $.trim(prescription) == "") {
      errorContainer.MaterialSnackbar.showSnackbar(
          {message : "Name, Age/Sex or Prescription cannot be empty!"});
      return
    }

    $.post('/api/genpdf',
           {name : name, age_sex : age_sex, prescription : prescription},
           function(data, status) {
             if (data.Err != '') {
               console.log(data.Err);
               errorContainer.MaterialSnackbar.showSnackbar(
                   {message : data.Err});
               return;
             }
             currPrescription = data.Data;
             window.open(data.Data, "_blank");
           });
  });

  // Clear Fields button handler.
  document.querySelector('#clear').addEventListener('click', function() {
    $('#name').val("");
    $('#age_sex').val("");
    $('#prescription').val("");
    currPrescription = "";
  });

  // Print PDF button handler.
  document.querySelector('#print').addEventListener('click', function() {
    if (currPrescription == "") {
      errorContainer.MaterialSnackbar.showSnackbar(
          {message : "Nothing to print!"});
      return;
    }
    $.post('/api/print', {file : currPrescription}, function(data, status) {
      if (data.Err != '') {
        console.log(data.Err);
        errorContainer.MaterialSnackbar.showSnackbar({message : data.Err});
        return;
      }
      errorContainer.MaterialSnackbar.showSnackbar({message : data.Data});
    });
  });
});
