<div class="right-col" role="main">
  <div>
    <div class="clearfix"></div>

    <!-- flash message -->
    {{ if .flash.notice }}
    <div class="alert alert-success alert-dismissible fade in" role="alert">
      <button type="button" class="close" data-dismiss="alert" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
      {{ .flash.notice }}
    </div>
    {{ end }}

    {{ if .flash.error }}
    <div class="alert alert-danger alert-dismissible fade in" role="alert">
      <button type="button" class="close" data-dismiss="alert" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
      {{ .flash.error }}
    </div>
    {{ end }}

    {{ if .flash.warning }}
    <div class="alert alert-warning alert-dismissible fade in" role="alert">
      <button type="button" class="close" data-dismiss="alert" aria-label="Close">
        <span aria-hidden="true">&times;</span>
      </button>
      {{ .flash.warning }}
    </div>
    {{ end }}
    <!-- /flash message -->

    <div class="clearfix"></div>

    {{ if .User }}
    <div class="x_panel">
      <div class="x_title">
        <h2>Edit User</h2>
      </div>
      <div class="x_content">
        <form method="POST" id="demo-form2" class="form-horizontal">
          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Email</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <p class="form-control-static" style="display: inline-block;">
                {{ .User.Email }}
                <!--{{ if eq .User.EmailVerified false}} (Unverified)
                <div class="col-md-5 pull-right">
                  <a href="#" class="btn btn-default"><span class="glyphicon glyphicon-envelope"></span> Send Verification Email</a>
                </div>
                {{ end }}-->
              </p>
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">New Password</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input type="password" id="password" name="password" class="form-control" data-parsley-trigger="change" data-parsley-pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}" data-parsley-pattern-message="Your password must contain a minimum of 8 characters and at least 1 upper case letter, 1 lower case letter and 1 number.">
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Confirm New Password</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
               <input type="password" id="confirm_password" name="confirm_password" class="form-control" data-parsley-equalto="#password" data-parsley-trigger="change" data-parsley-equalto-message="Passwords do not match.">
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Name</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input id="alias" name="alias" type="text" class="form-control" data-parsley-maxlength="255" data-parsley-trigger="change" placeholder="" value="{{ .User.Alias }}">
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Role</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              {{ if eq .EnableAdvancedOptions false }}
              <select id="roleUpdate" name="roleUpdate" class="form-control" value="{{ .User.Role }}" disabled>
              {{ else }}
              <select id="roleUpdate" name="roleUpdate" class="form-control" value="{{ .User.Role }}">
              {{ end }}
                <option value="admin">Administrator</option>
                <option value="user">User</option>
              </select>
            </div>
          </div>
          {{if .flash.error}}
          <div class="form-group">
            <div class="col-md-6 col-md-offset-4">
              <font color="red">{{.flash.error}}</font>
            </div>
          </div>
          {{end}}
          <!-- <div class="ln_solid"></div> -->
          <br />
          <div class="form-group pull-right">
            <div class="col-md-12">
              <button type="submit" class="btn btn-primary" data-disable-with="Save Changes">Save Changes</button>
              <a href="/users/{{ .User.Id }}" class="btn btn-default">Cancel</a>
            </div>
          </div>
        </form>
      </div>
    </div>

    <div class="clearfix"></div>

    {{ if .EnableAdvancedOptions }}
    {{ if eq .EnableAdvancedOptions true}}
    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Advanced Settings</h2>
            <div class="clearfix"></div>
          </div>

          <div class="x_content">
            {{ if eq .User.Enabled true }}
            <p>You have to DISABLE USER before DELETING USER.</p>
            <button onclick="launchModal.call(this, 'Disable User', 'userDisableTemplate', 'disable');" type="button" class="btn btn-warning" data-profile-id="{{ .User.Id }}"><i class="ion ion-power"></i> Disable User</button>
            {{ else }}
            <button onclick="launchModal.call(this, 'Enable User', 'userEnableTemplate', 'enable');" type="button" class="btn btn-success" data-profile-id="{{ .User.Id }}"><i class="ion ion-power"></i> Enable User</button>
            {{ end }}
            {{ if eq .User.Enabled true }}
            <button onclick="launchModal.call(this, 'Delete User', 'userDeleteTemplate', 'delete');" type="button" class="btn btn-danger" data-profile-id="{{ .User.Id }}" disabled="disabled"><i class="ion ion-trash-a"></i> Delete User</button>
            {{ else }}
            <button onclick="launchModal.call(this, 'Delete User', 'userDeleteTemplate', 'delete');" type="button" class="btn btn-danger" data-profile-id="{{ .User.Id }}"><i class="ion ion-trash-a"></i> Delete User</button>
            {{ end }}
          </div>
          {{ end }}
        </div>
      </div>
    </div>
    {{ end }}
    {{ end }}
  </div>
</div>

<!-- Dialogue -->
<!-- Launcher -->
<div id="modalLauncher" class="modal fade" tabindex="-1" role="dialog" aria-hidden="true"></div>
<!-- Template -->
<div id="modalTemplate" style="display: none;">
  <div class="modal-dialog">
    <div class="modal-content">
      <form method="POST" id="action_form" class="form-horizontal">
        <div class="modal-header">
          <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span>
            </button>
          <h4 class="modal-title" id="modal-title">ACTION_TITLE</h4>
        </div>
        <div class="modal-body" id="modal-body">
          ACTION_CONTENT
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
          <span id="modal-action"></span>
        </div>
      </form>
    </div>
  </div>
  <div id="user_id_refer" style="display: none;">{{ .User.Id }}</div>
</div>
<!-- /Dialogue -->
<!-- Dialogue content: disable user-->
<div id="userDisableTemplate" style="display: none">
  <h4>Are you sure you want to do this?</h4>
  <div class="panel panel-warning">
    <div class="panel-heading">By performing the <b>Disable User</b> action:</div>
    <ul class="list-group">
      <li class="list-group-item">Prevent the User from logging into the Subspace Management Console.</li>
      <li class="list-group-item">You will immediately disconnect all the VPN sessions associated with User.</li>
      <li class="list-group-item">Profiles associated with this User will be disabled and unusable.</li>
    </ul>
  </div>
</div>
<!-- /Dialogue content: disable user -->
<!-- Dialogue content: enable user-->
<div id="userEnableTemplate" style="display: none">
  <h4>Are you sure you want to do this?</h4>
  <div class="panel panel-success">
    <div class="panel-heading">By performing the <b>Enable User</b> action:</div>
    <ul class="list-group">
      <li class="list-group-item">If the User's Role is an Administrator, the User will be able to log into the Subspace Management Console.</li>
      <li class="list-group-item">The User's associated Profiles will not be renabled as a result. Each Profile will need to be renabled manually for security purposes.</li>
    </ul>
  </div>
</div>
<!-- /Dialogue content: enable user -->
<!-- Dialogue content: delete user -->
<div id="userDeleteTemplate" style="display: none;">
  <h4>Are you sure you want to do this?</h4>
  <div class="panel panel-danger">
    <div class="panel-heading">By performing the <b>Delete User</b> action:</div>
    <ul class="list-group">
      <li class="list-group-item">You will immediately disconnect all the VPN sessions associated with this User.</li>
      <li class="list-group-item">This User will be deleted and all associated Profiles will be deleted.</li>
    </ul>
  </div>
  To complete this action, type in the User's email: <b><span class="deleteEmail">{{ .User.Email }}</span></b>
  <input id="confirmDeleteEmail" name="confirmDeleteEmail" type="email" class="form-control" value="" placeholder="">
  <div class="clearfix"></div>
</div>
<!-- /Dialogue content: delete user -->
