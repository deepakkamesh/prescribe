

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

  // Generate PDF button handler.
  document.querySelector('#generate_pdf').addEventListener('click', function() {
    name = $('#name').val();
    age_sex = $('#age_sex').val();
    prescription = $('#prescription').val();

    $.post('/api/genpdf',
           {name : name, age_sex : age_sex, prescription : prescription},
           function(data, status) {
             if (data.Err != '') {
               console.log(data.Err);
               errorContainer.MaterialSnackbar.showSnackbar(
                   {message : data.Err});
               return;
             }
             window.open(data.Data, "_blank");
           });
  });

  // Clear Fields button handler.
  document.querySelector('#clear').addEventListener('click', function() {
    $('#name').val("");
    $('#age_sex').val("");
    $('#prescription').val("");
  });

  // Print PDF button handler.
  document.querySelector('#print').addEventListener('click', function() {
    console.log("dd");
    errorContainer.MaterialSnackbar.showSnackbar({message : "sdsd"});
  });
});
