

$(document).ready(function() {
  var errorContainer = document.querySelector('#error-popup');
  // Generate PDF.
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
               return
             }
             window.open(data.Data, "_blank");
           });
  });

  // Clear Fields.
  document.querySelector('#clear').addEventListener('click', function() {
    $('#name').val("");
    $('#age_sex').val("");
    $('#prescription').val("");
  });

  // Print PDF.
  document.querySelector('#print').addEventListener('click', function() {
    console.log("dd");
    errorContainer.MaterialSnackbar.showSnackbar({message : "sdsd"});
  });
});
