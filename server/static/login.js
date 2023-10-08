$(document).ready(function() {
  addFormValidation();

  $('#login').submit(function(e) {
    const isValid = $('#login').form('validate form');
    if (!isValid) {
      return false;
    }
    e.preventDefault();
    ajaxLogin();
  });
});

function ajaxLogin() {
  jsonData = {
    email: $('#email').val(),
    password: $('#password').val()
  };

  $.ajax({
    type: 'POST',
    url: 'http://127.0.0.1:9001/api/login',
    data: JSON.stringify(jsonData),
    contentType: 'application/json',
    encode: true,
    success: function (data, textStatus, xhr) {
      document.cookie = "chatToken=" + data.accessToken;
      window.location.replace("/chat");
    },
    error: function(jqXHR, textStatus, errorThrown) {
      console.log(jqXHR, textStatus, errorThrown);
    },
  });
}//ajaxLogin

function addFormValidation() {
  $('#login')
    .form({
      on: 'submit',
      fields: {
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
