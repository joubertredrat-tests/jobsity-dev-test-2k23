$(document).ready(function() {
  $('.ui.dropdown').dropdown();
  const token = parseJwt(getCookie('chatToken'));
  $('#userName').html(token.userName);
  ajaxMessages();
  setInterval(ajaxMessages, 1000);
});//$(document).ready

function parseJwt(token) {
  var base64Url = token.split('.')[1];
  var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
  var jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));

  return JSON.parse(jsonPayload);
}//parseJwt

function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(';').shift();
}//getCookie

function ajaxMessages() {
  $.ajax({
    type: 'GET',
    url: 'http://127.0.0.1:9001/api/messages',
    headers: {'Authorization': 'Bearer ' + getCookie('chatToken')},
    encode: true,
    success: function (data, textStatus, xhr) {
      buildChat(data);
    },
    statusCode: {
      500: function() {
        Swal.fire({
          icon: 'error',
          title: 'Error!',
          text: 'Anything wrong is not right, try again later'
        });
      }
    }
  });
}//ajaxMessages

function buildChat(messages) {
  messages.reverse().forEach(function(message) {
    if (messageExists(message.id)) {
      return;
    }

    $('#chat').append(getMessageHtml(message));
  });

}//buildChat

function messageExists(id) {
  return $('#message-' + id).length > 0;
}//messageExists

function getMessageHtml(message) {
  return `<div class="ui ${getMessageHtmlRelatedClass(message.userEmail)} message" id="message-${message.id}">
    <div class="header">
      ${message.userName} - ${formatDate(message.datetime)}
    </div>
    <p>${message.messageText}</p>
  </div>`;
}//getMessageHtml

function getMessageHtmlRelatedClass(userEmail) {
  const token = parseJwt(getCookie('chatToken'));
  return token.userEmail == userEmail ? 'olive' : 'grey';
}//getMessageHtmlRelatedClass

function formatDate(date) {
  if (date == null){
    return '';
  }

  const d = new Date(date);
  return (
    [
      padTo2Digits(d.getMonth() + 1),
      padTo2Digits(d.getDate()),
      d.getFullYear(),
    ].join('/') +
    ' ' +
    [
        padTo2Digits(d.getHours()),
        padTo2Digits(d.getMinutes()),
    ].join(':')
  );
}//formatDate

function padTo2Digits(num) {
  return num.toString().padStart(2, '0');
}//padTo2Digits
