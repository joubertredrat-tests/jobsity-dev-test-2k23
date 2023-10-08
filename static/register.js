$(document).ready(function() {
  addFormValidation();

  $('#register').submit(function(e) {
    const isValid = $('#register').form('validate form');
    if (!isValid) {
      return false;
    }
    e.preventDefault();
    ajaxRegister();
  });
});//$(document).ready

function ajaxRegister() {
  jsonData = {
    name: $('#name').val(),
    email: $('#email').val(),
    password: $('#password').val()
  };

  $.ajax({
    type: 'POST',
    url: 'http://127.0.0.1:9001/api/register',
    data: JSON.stringify(jsonData),
    contentType: 'application/json',
    encode: true,
    success: function (data, textStatus, xhr) {
      window.location.replace("/login");
    },
    statusCode: {
      400: function() {
        Swal.fire({
          icon: 'error',
          title: 'Error!',
          text: 'Some fields are filled with wrong data, try again'
        });
      },
      422: function() {
        Swal.fire({
          icon: 'error',
          title: 'Error!',
          text: 'E-mail filled already registered, try login'
        });
      },
      500: function() {
        Swal.fire({
          icon: 'error',
          title: 'Error!',
          text: 'Anything wrong is not right, try again later'
        });
      }
    }
  });
}//ajaxLogin

function addFormValidation() {
  $('#register')
    .form({
      on: 'submit',
      fields: {
        name: {
          identifier: 'name',
          rules: [
            {
              type: 'empty',
              prompt: 'Fill your name'
            }
          ]
        },
        email: {
          identifier: 'email',
          rules: [
            {
              type: 'email',
              prompt: 'Fill correct e-mail'
            }
          ]
        },
        password: {
          identifier: 'password',
          rules: [
            {
              type   : 'minLength[8]',
              prompt : 'Your password must be at least {ruleValue} characters'
            }
          ]
        }
      }
    })
  ;
}//addFormValidation
