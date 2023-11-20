<div>
  <div class="row">
    <div class="col-md-6 col-sm-6 col-xs-12 subspace-create-admin-form-border">
      <div class="x_panel split">
        <div class="x_title subspace-h1-title-visible">
          <h1>Create Administrator Account</h1>
        </div>

        <div class="x_content">
          <form method="POST" action="" id="create-admin-form" data-parsley-validate="" class="form-horizontal">
            <div style="padding-bottom: 20px">This account will be used to manage your Subspace service.</div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="alias">
                  Name
                </label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="text" id="alias" name="alias" class="form-control" data-parsley-maxlength="255" data-parsley-trigger="change" placeholder="Name" value="{{ .Form.Alias }}" />
              </div>
            </div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="email">
                  Email Address
                  <span class="required">*</span>
                </label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="email" id="email" name="email" class="form-control" placeholder="Email Address" data-parsley-maxlength="255" data-parsley-required="true" data-parsley-trigger="change" value="{{ .Form.Email }}" required />
              </div>
            </div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="password">
                    Password
                    <span class="required">*</span>
                </label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="password" id="password" name="password" data-parsley-required="true" class="form-control" placeholder="Password" data-parsley-trigger="change" data-parsley-errors-container="#password_error" data-parsley-pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}" data-parsley-pattern-message="Your password must contain a minimum of 8 characters and at least 1 upper case letter, 1 lower case letter and 1 number." required />
              </div>
            </div>
            <div class="form-group">
              <label class="col-md-12 col-sm-12 col-xs-12" for="confirm_password">
                  Confirm Password
                  <span class="required">*</span>
                </label>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="password" id="confirm_password" name="confirm_password" data-parsley-required="true" data-parsley-equalto="#password" data-parsley-trigger="change" class="form-control" data-parsley-equalto-message="Passwords do not match." placeholder="Type Your Password Again"
                  required>
              </div>
            </div>
            {{if .flash.error}}
            <div>
              <font color="red">{{.flash.error}}</font>
            </div>
            {{end}}
            <div class="setup-action-Btn">
              <button type="submit" class="btn btn-block btn-success" data-disable-with="Create">Create</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <div class="col-md-6 col-sm-6 col-xs-12">
      <div class="x_panel split">
        <div class="x_title subspace-h1-title-visible">
          <h1>Restore from Backup</h1>
        </div>
        <div class="x_content">
          <form id="restore" method="POST" action="" class="form-horizontal" enctype="multipart/form-data">
            <div class="form-group" style="margin-bottom:15px;">
              <div class="col-md-12 col-sm-12 col-xs-12" style="padding-bottom:32px">
                <p>
                  Upload a Subspace configuration file to restore settings.
                </p>
              </div>
              <div class="col-md-12 col-sm-12 col-xs-12">
                <input type="file" name="upload_config" id="upload_config" class="" title="Choose file..." />
              </div>
            </div>
            <div class="form-group">
              <div class="col-md-12 col-sm-12 col-xs-12">
                <div id="finish_restore" class="alert alert-success alert-dismissible subspace-alert-style fade in" role="alert" style="display:none;">
                  <button type="button" class="close" data-dismiss="alert" aria-label="Close">
                    <span aria-hidden="true">×</span>
                  </button>
                  <span id="restore_completed_msg"></span>
                </div>
                <div id="failed_restore" class="alert alert-danger alert-dismissible subspace-alert-style fade in" role="alert" style="display:none;">
                  <button type="button" class="close" data-dismiss="alert" aria-label="Close"></button>
                  <span id="restore_failed_msg">
                    Restore unsuccessful. Please choose another file to try again or contact support.
                  </span>
                </div>
                <button id="restore_again" type="button" form="restore" class="btn btn-primary" style="display:none;">
                  <i class="ion ion-archive" style="font-size:16px; padding-right:8px"></i>
                  Restore again
              </button>
              </div>
            </div>
            <div class="form-group">
              <div class="col-md-12">
                <button id="start_restore" type="button" form="restore" class="btn btn-primary">
                  <i class="ion ion-archive" style="font-size:16px; padding-right:8px"></i>
                  Restore
                </button>
                <button id="restoring" class="btn btn-danger" type="button" style="display:none;" disabled>
                  <i class="fa fa-spinner fa-pulse fa-lg fa-fw"></i>
                  Restoring...
                </button>
                <p id="restore_message" style="color:red;margin-top:10px;"></p>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Button to close the overlay navigation -->
<!-- Overlay content -->
<div id="myNav" class="overlay">
  <div class="overlay-content">
    <a href="#"><i class="fa fa-spinner fa-pulse fa-lg fa-fw"></i>  Restoring...</a>
    <br>
    <p>Please wait while the your settings are restored. The whole process may take a few minutes.</p>
    <p>You will be redirected to the log in page after the restoration process is complete. Thank you for your patience.</p>
  </div>
</div>
