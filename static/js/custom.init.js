const restoreConfigFile = document.getElementById("upload_config");
const start_restore = document.getElementById("start_restore");
const restore_again = document.getElementById("restore_again");
var ask, askTimeout, restoreFileURI;

var askStatus = (option) => {
  fetch('/init/clean_install?backup_restore_option=' + option, {
      method: 'get',
      credentials: 'same-origin'
    })
    .then((response) => {
      return response.json();
    })
    .then((json) => {
      console.log("json: ", json);
      result = JSON.stringify(json);

      if (option == "restore") {
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
        // location.reload();
        window.location = "/sign_in";
      }, 10000);
    })
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
  // start_backup.disabled = false;
};

const startRestore = (event) => {
  const file = restoreConfigFile.files[0];
  var option = "restore"
  if (file !== undefined) {
    // start_backup.disabled = true;
    openNav();
    $("#restore_message").text("");
    var formData = new FormData();
    formData.append("backup_restore_option", "restore");
    formData.append("upload_config", file);
    fetch('/init/clean_install', {
        method: 'POST',
        body: formData,
        credentials: 'same-origin'
      })
      .then(response => {
        return response.json();
      })
      .then(success => {
        // console.log("response: ", success);
        if (success == 'false') {
          closeNav();
          $("#restore_message").text("Restore process is unable to start.");
        } else {
          restoreFileURI = success;
          $("#start_restore").hide();
          $("#restoring").fadeIn();
          ask = setInterval(() => {
            askStatus(option)
          }, 2000);
        }
        // askTimeout = setTimeout(stopAskingRestoreStatus, 10000);
      }).catch(error => {
        console.log("error: ", error);
      });
  } else {
    $("#restore_message").text("Please select a backup file to upload.");
  }
};

const restoreAgain = () => {
  $("#upload_config").show();
  $("#restore_again").hide();
  $("#finish_restore").hide();
  $("#failed_restore").hide();
  $("#start_restore").fadeIn();
}

// Add a listener on input
if (start_restore !== null) {
  start_restore.addEventListener("click", startRestore, false);
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
