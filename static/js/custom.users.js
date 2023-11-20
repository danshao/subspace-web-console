function launchModal(modalTitle, modalContentRef, modalAction) {
  var HTTP_PATCH  = '<input type="hidden" name="_method" value="patch" />',
      HTTP_PUT    = '<input type="hidden" name="_method" value="put" />',
      HTTP_DELETE = '<input type="hidden" name="_method" value="delete" />';

  var userID = $("#modalTemplate #user_id_refer").html();
  var profileID = $(this).data("profile-id");
  var sessionName = $(this).data("session-id");
  var action = modalAction;

  $("#modalTemplate #modal-title").html(modalTitle);
  $("#modalTemplate #modal-body").html($("#" + modalContentRef).html());

  if (modalContentRef.indexOf("user") != -1) {
    var uri = "/users/" + userID + "/" + action;
    $("#modalTemplate #action_form").attr("action", uri);

    if (modalContentRef == "userDeleteTemplate") {
      $("#modalTemplate #modal-action").html('<button id="deleteUserConfirmBtn" type="submit" class="btn btn-danger" form="action_form" value="Submit" disabled="disabled">Delete</button>');
    } else if (modalContentRef == "userDisableTemplate") {
      $("#modalTemplate #modal-action").html(HTTP_PATCH + '<button id="modal-action" type="submit" class="btn btn-warning" form="action_form" value="Submit">Disable</button>');
    } else {
      $("#modalTemplate #modal-action").html(HTTP_PATCH + '<button id="modal-action" type="submit" class="btn btn-success" form="action_form" value="Submit">Enable</button>');
    }
  } else {
    if (modalContentRef == "profileEditTemplate") {
      description = $(this).parents("tr").prev().find("td").eq(0).html();
      $("#" + modalContentRef + " #edit-description").attr("value", description);
    }

    var uri = "/users/" + userID + "/" + action + "/" + profileID

    switch (modalContentRef) {
    case "profileDeleteTemplate":
      $("#modalTemplate #modal-action").html('<a id="modal-action" href="' + uri + '" class="btn btn-danger" role="button">Delete</a>');
      break;

    case "profileEditTemplate":
      $("#modalTemplate #modal-action").html('<button type="submit" class="btn btn-primary" form="action_form" value="submit">Save Changes</button>');
      $("#modalTemplate #action_form").attr("action", uri)
      break;

    case "profileDisableTemplate":
      uri = ["", "users", userID, "profiles", profileID, action].join("/")
      $("#modalTemplate #action_form").attr("action", uri)
      $("#modalTemplate #modal-action").html(HTTP_PATCH + '<button id="modal-action" type="submit" class="btn btn-warning" form="action_form" value="Submit">Disable</button>');
      break;

    case "profileDisconnecTemplate":
      uri = "/users/" + userID + "/" + action + "/" + sessionName
      $("#modalTemplate #modal-action").html('<a id="modal-action" href="' + uri + '" class="btn btn-danger" role="button">Disconnect</a>');
      break;

    case "profileEnableTemplate":
      uri = ["", "users", userID, "profiles", profileID, action].join("/")
      $("#modalTemplate #action_form").attr("action", uri)
      $("#modalTemplate #modal-action").html(HTTP_PATCH + '<button id="modal-action" type="submit" class="btn btn-success" form="action_form" value="Submit">Enable</button>');
      break;

    default:
      // Error condition
      $("#modalTemplate #modal-action").html("");
    }
  }

  $("#modalTemplate #modal-title").html(modalTitle);
  $("#modalTemplate #modal-body").html($("#" + modalContentRef).html());
  $("#modalLauncher").html($("#modalTemplate").html()).modal('show');

  // Disable Delete button if email does not match
  $(function() {
    $('#confirmDeleteEmail').keyup(function() {
      if ($(this).val() == $('.deleteEmail').html()) {
        $('#deleteUserConfirmBtn').prop("disabled", false);
      } else {
        $('#deleteUserConfirmBtn').prop("disabled", true);
      }
    });
  });
}

// $('#modalLauncher').on('show.bs.modal', function (e) {
//     $(this).find('.modal-body').html();
//     $(this).html($("#modalTemplate").html());
//     console.log(e);
// })

// Edit User - Set user's role as selected option
$(function() {
  $('#roleUpdate').val($('#roleUpdate').attr('value')).prop('selected', true);
});

// Parsley
$(function() {
  $('#create-user-form').parsley({
    excluded: '[disabled]',
  });

  $('#demo-form2').parsley();
});

// Create User - Password Input fields are disabled/enabled based on if autogenerate password is checked or unchecked.
$(function() {
  $('#autogenPassword').on('click', function(event) {
    if ($('#autogenPassword').prop('checked') === true) {
      $('#password').prop('disabled', true);
      $('#confirm_password').prop('disabled', true);

      $('#password').prop('required', false);
      $('#confirm_password').prop('required', false);

      $('#password').parsley().reset();
      $('#confirm_password').parsley().reset();

      $('#password').val("");
      $('#confirm_password').val("");
    } else {
      $('#password').prop('disabled', false);
      $('#confirm_password').prop('disabled', false);

      $('#password').prop('required', true);
      $('#confirm_password').prop('required', true);

      $('#password').parsley();
      $('#confirm_password').parsley();
    }
  })
});
