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

    <div class="x_panel">
      <div class="x_title">
        <h2>Create New User</h2>
      </div>
      <div class="x_content">
        <form method="POST" id="create-user-form" class="form-horizontal">

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12"><span class="required">*</span>Email Address</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input id="email" name="email" type="email" class="form-control" data-parsley-maxlength="255" data-parsley-required="true" data-parsley-trigger="change" required>
            </div>
          </div>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12"><span class="required">*</span>Password</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input type="password" id="password" name="password" data-parsley-required="true" class="form-control" data-parsley-trigger="change" data-parsley-pattern="(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{6,}" data-parsley-pattern-message="Your password must contain a minimum of 8 characters and at least 1 upper case letter, 1 lower case letter and 1 number." disabled>
            </div>
          </div>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12"><span class="required">*</span>Confirm Password</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input type="password" id="confirm_password" name="confirm_password" class="form-control" data-parsley-required="true" data-parsley-equalto="#password" data-parsley-trigger="change" data-parsley-equalto-message="Passwords do not match." disabled>
              <div class="checkbox">
                <label class="label-cr">
                  <input id="autogenPassword" name="autogenPassword" type="checkbox" checked="checked" />
                  <span class="labelbox"><i class="ico"></i>Generate Random Password</span>
                </label>
              </div>
            </div>
          </div>

          <br>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Name</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <input id="alias" name="alias" type="text" data-parsley-maxlength="255" data-parsley-trigger="change" class="form-control" placeholder="">
            </div>
          </div>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Role</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <select id="role" name="role" class="form-control" value="user">
                <option value="admin"> Administrator </option>
                <option value="user"> User </option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">VPN</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <div class="checkbox">
                <label class="label-cr">
                  <input id="createVPNProfile" name="createVPNProfile" type="checkbox" value="1" checked="checked">
                  <span class="labelbox"><i class="ico"></i>Automatically create a VPN Profile for this user</span>
                </label>
              </div>
            </div>
          </div>

          <div class="form-group">
            <label class="control-label col-md-3 col-sm-3 col-xs-12">Delivery Method</label>
            <div class="col-md-6 col-sm-6 col-xs-12">
              <div class="checkbox">
                <label class="label-cr">
                  {{ if eq .CanSendEmail false }}
                  <input id="delivery" name="delivery" type="checkbox" value="0" disabled>
                  {{ else }}
                  <input id="delivery" name="delivery" type="checkbox" value="1" checked="checked">
                  {{ end }}
                  <span class="labelbox">
                    <i class="ico"></i>Send instructions to user via email
                    {{ if eq .CanSendEmail false }}
                    <i class="ion ion-information-circled" data-toggle="tooltip" data-placement="right" title="" data-original-title="Go to 'Settings' > 'Mail' to enable this feature" style="padding-left:5px"></i>
                    {{ end }}
                  </span>
                </label>
              </div>
            </div>
          </div>

          <div class="ln_solid"></div>

          <div class="form-group pull-right">
            <div class="col-md-12">
              <button type="submit" class="btn btn-primary" data-disable-with="Create User">Create User</button>
              <a href="/users" class="btn btn-default">Cancel</a>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
