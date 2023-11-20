function launchModal(modalTitle, modalContentRef, modalAction) {

  var uri = "/settings/data_update/" + modalAction
  var description = $(this).parents("div").children("span.refer").html();

  $("#modalTemplate #modal-title").html(modalTitle);
  $("#modalTemplate #action_form").attr("action", uri);
  $("#" + modalContentRef + " .edit-description").attr("value", description);

  if (modalContentRef == "hostnameEditTemplate") {
    $("#modalTemplate #modal-action").html('<button id="hostnameEditBtn" type="submit" class="btn btn-primary" form="action_form" value="submit" disabled>Save Changes</button>');
  } else if (modalContentRef == "presharedKeyEditTemplate") {
    $("#modalTemplate #modal-action").html('<button id="EditBtn" type="submit" class="btn btn-primary" form="action_form" value="submit">Save Changes</button>');
  } else {
    $("#modalTemplate #modal-action").html('<button type="submit" class="btn btn-primary" form="action_form" value="submit">Re-generate</button>');
  }
  $("#modalTemplate #modal-body").html($("#" + modalContentRef).html());
  $("#modalLauncher").html($("#modalTemplate").html()).modal('show');

  if (modalContentRef == "hostnameEditTemplate") {
    if ($("#action_form input[name$='hostname']").val() == "") {
      $("#hostnameEditBtn").prop("disabled", false);
    }
    $("#action_form input[name$='hostname']").on("keyup", () => {
      var hostname = $("#action_form input[name$='hostname']").val();
      if (hostname == "" || /^((([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.)+([a-zA-Z]{2,61})$/.test(hostname)) {
        $("#hostnameEditBtn").prop("disabled", false);
        $("#action_form .error-msg").html("");
      } else {
        $("#hostnameEditBtn").prop("disabled", true);
        $("#action_form .error-msg").html("incorrect domain name format");
      }
    });
  }

  if (modalContentRef == "presharedKeyEditTemplate") {
    if ($("#action_form input[name$='presharedkey']").val() == "") {
      $("#EditBtn").prop("disabled", false);
    }
    $("#action_form input[name$='presharedkey']").on("keyup", () => {
      var presharedkey = $("#action_form input[name$='presharedkey']").val();
      if (/^[a-zA-Z0-9]{1,9}$/.test(presharedkey)) {
        $("#EditBtn").prop("disabled", false);
        $("#action_form .error-msg").html("");
      } else {
        $("#EditBtn").prop("disabled", true);
        $("#action_form .error-msg").html("Pre-shared key must be alphanumeric and 9 characters or less.");
      }
    });
  }

}

// Mail Settings - Authentication options enabled/disabled based on if option is checked.
$(function() {
  $('#authentication').on('click', function(event) {
    if ($('#authentication').prop('checked') === true) {
      $('#username').prop('disabled', false);
      $('#password').prop('disabled', false);

      $('#username').prop('required', true);
      $('#password').prop('required', true);
    } else {
      $('#username').prop('disabled', true);
      $('#password').prop('disabled', true);

      $('#username').prop('required', false);
      $('#password').prop('required', false);
    }
  })
});

var verifyDNS = (hostname) => {
  $("#dns-verify").removeClass("btn-success btn-danger").addClass("btn-default");
  $("#dns-verify .fa-check-circle-o").hide();
  $("#dns-verify .fa-times-circle-o").hide();
  $("#dns-verify .fa-spinner").fadeIn("fast");
  var formData = new FormData();
  formData.append("hostname", hostname);
  fetch('/settings/dns_verify', {
      method: 'post',
      body: formData
    })
    .then((response) => {
      return response.json();
    })
    .then((json) => {
      var res = JSON.stringify(json);
      if (res == "true") {
        $("#host-verify-message").hide().html("");
        $("#dns-verify .fa-spinner").hide();
        $("#dns-verify .fa-times-circle-o").hide();
        $("#dns-verify .fa-check-circle-o").fadeIn("fast");
        $("#dns-verify").removeClass("btn-default").addClass("btn-success");
      } else {
        $("#host-verify-message").html("Hostname does not match the IP of Subspace").show();
        $("#dns-verify .fa-spinner").hide();
        $("#dns-verify .fa-check-circle-o").hide();
        $("#dns-verify .fa-times-circle-o").fadeIn("fast");
        $("#dns-verify").removeClass("btn-default").addClass("btn-danger");
      }
      // console.log("response: ", json);
    })
    .catch((reason) => {
      console.log("verfify DNS error");
      console.log(reason);
    })
}

const restoreConfigFile = document.getElementById("upload_config");
const start_restore = document.getElementById("start_restore");
const start_backup = document.getElementById("start_backup");
const backup_again = document.getElementById("backup_again");
const restore_again = document.getElementById("restore_again");
var ask, askTimeout, backupFileURI, restoreFileURI;

var askStatus = (option) => {
  fetch('/settings/ajax_response?backup_restore_option=' + option, {
      method: 'get',
      credentials: 'same-origin'
    })
    .then((response) => {
      return response.json();
    })
    .then((json) => {
      console.log("json: ", json);
      result = JSON.stringify(json);

      if (option == "backup") {
        switch (true) {
          case /running/.test(result):
            // do nothing
            break;
          case /failed/.test(result):
            stopAskingBackupStatus('failed');
            $("#backup_failed_msg").text(result);
            backupFileURI = '';
            break;
          default:
            stopAskingBackupStatus();
            fileArray = backupFileURI.split('/')
            $("#download_backup").attr("href", backupFileURI);
            var backupCompletedMsg = "Backup config: " + fileArray[fileArray.length-1] + " is ready to download, click the button below."
            $("#backup_completed_msg").text(backupCompletedMsg);
            backupFileURI = '';
            break;
        }
      } else if (option == "restore") {
        switch (true) {
          case /running/.test(result):
            // do nothing
            break;
          case /failed/.test(result):
            closeNav();
            stopAskingRestoreStatus('failed');
            // $("#restore_failed_msg").text(result);
            restoreFileURI = '';
            restoreCompletedMsg = '';
            break;
          case /succeeded/.test(result):
            stopAskingRestoreStatus();
            window.location = "/sign_in";
            break;
          default:
            closeNav();
            stopAskingRestoreStatus();
            restoreCompletedMsg = "Restore process using config: " + restoreFileURI + " is accomplished, proceed to use VPN service."
            $("#restore_completed_msg").text(restoreCompletedMsg);
            restoreFileURI = '';
            restoreCompletedMsg = '';
            break;
        }
      }
    })
    .catch((reason) => {
      //console.log("err: ", reason);
      clearInterval(ask);
      //TODO Too fast will redirect to Maintain page. Fix the ugly solution.
      setTimeout(function(){
        location.reload();
      }, 5000);
    })
};

var stopAskingBackupStatus = (status) => {
  clearInterval(ask);
  $("#backing_up").hide();
  if (status == 'failed') {
    $("#failed_backup").fadeIn();
  } else {
    $("#download_backup").fadeIn();
    $("#finish_backup").fadeIn();
  }
  $("#backup_again").fadeIn();
  start_restore.disabled = false;
};

var stopAskingRestoreStatus = (status) => {
  clearInterval(ask);
  closeNav();
  $("#restoring").hide();
  $("#upload_config").val("");
  $("#upload_config").parent().next('.file-input-name').remove()
  $("#upload_config").hide();
  if (status == 'failed') {
    $("#failed_restore").fadeIn();
  } else {
    $("#finish_restore").fadeIn();
  }
  // $("#restore_again").fadeIn();
  $("#upload_config").show();
  $("#start_restore").fadeIn();
  start_backup.disabled = false;
};

const startBackup = () => {
  start_restore.disabled = true;
  var formData = new FormData();
  var option = "backup";
  formData.append("backup_restore_option", "backup");
  fetch('/settings/ajax_response', {
      method: 'POST',
      body: formData,
      credentials: 'same-origin'
    })
    .then(response => {
      return response.json();
    })
    .then(success => {
      // console.log("response: ", success);
      backupFileURI = success;
      $("#start_backup").hide();
      $("#backing_up").fadeIn();
      ask = setInterval(() => {
        askStatus(option)
      }, 2000);
      // askTimeout = setTimeout(stopAskingBackupStatus, 10000);
    }).catch(error => {
      console.log("error: ", error);
    });
};

const startRestore = (event) => {
  const file = restoreConfigFile.files[0];
  var option = "restore"
  if (file !== undefined) {
    start_backup.disabled = true;
    openNav();
    $("#restore_message").text("");
    var formData = new FormData();
    formData.append("backup_restore_option", "restore");
    formData.append("upload_config", file);
    fetch('/settings/ajax_response', {
        method: 'POST',
        body: formData,
        credentials: 'same-origin'
      })
      .then(response => {
	console.log("response:", response)
        console.log("response headers:", response.headers);
        switch (response.status) {
        case 200:
          success = response.json();
          console.log("response success: ", success);
          if ("false" != success) {
            restoreFileURI = success;
            $("#start_restore").hide();
            $("#restoring").fadeIn();
            ask = setInterval(() => {
              askStatus(option)
            }, 2000);
          } else {
            closeNav();
            $("#restore_message").text("Restore process is unable to start.");
          }
          break;
        default:
          // If any error then redirect to root.
          window.location.replace("/");
          break;
        }
      })
      .catch(error => {
        console.log("error: ", error);
      });
  } else {
    $("#restore_message").text("Please select a backup file to upload.");
  }
};

const backupAgain = () => {
  $("#backup_again").hide();
  $("#finish_backup").hide();
  $("#download_backup").hide();
  $("#start_backup").fadeIn();
}
const restoreAgain = () => {
  $("#upload_config").show();
  $("#restore_again").hide();
  $("#finish_restore").hide();
  $("#failed_restore").hide();
  $("#start_restore").fadeIn();
}

// Add a listener on input
if (start_backup !== null) {
  start_backup.addEventListener("click", startBackup, false);
}
if (start_restore !== null) {
  start_restore.addEventListener("click", startRestore, false);
}
if (backup_again !== null) {
  backup_again.addEventListener("click", backupAgain, false);
}
if (restore_again !== null) {
  restore_again.addEventListener("click", restoreAgain, false);
}

// Screen Overlay Effect
// Open when someone clicks on the span element
function openNav() {
    document.getElementById("myNav").style.height = "100%";
}
// Close when someone clicks on the "x" symbol inside the overlay
function closeNav() {
    document.getElementById("myNav").style.height = "0%";
}
