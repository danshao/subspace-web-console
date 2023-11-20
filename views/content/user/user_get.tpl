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

    <!-- User Data -->
    <div class="page-title">
      <div class="title_left" style="width: 100%;">
        <h2>
          User Information
          {{ if .User }}
          <a href="/users/{{ .User.Id }}/edit" class="btn btn-sm btn-primary vertical-align-top pull-right"><i class="ion ion-edit"></i> Edit User</a>
          {{ end }}
        </h2>
      </div>
    </div>

    <div class="x_panel">
      <div class="x_content">
        {{ if .User }}
        <dl class="dl-horizontal">
          <dt>Email:</dt>
          <dd>{{ .User.Email }}</dd>
          <dt>Alias:</dt>
          <dd>{{ .User.Alias }}</dd>
          <dt>Role:</dt>
          <dd>
            {{ if eq .User.Role "admin"}}
            Administrator
            {{ else }}
            User
            {{ end }}
          </dd>
          <dt>Status:</dt>
          <dd>
            {{ if eq .User.Enabled true}}
            Enabled
            {{ else }}
            Disabled
            {{ end }}
          </dd>
          <dt>Date Created:</dt>
          <dd>{{ localTimeFmt .User.CreatedDate }}</dd>
          <dt>Date Modified:</dt>
          <dd>{{ localTimeFmt .User.UpdatedDate }}</dd>
        </dl>
        {{ end }}
      </div>
    </div>
    <!-- /User Data -->

    <hr>

    <div class="clearfix"></div>

    <!-- VPN Profile -->
    <div class="page-title">
      <div class="title_left" style="width: 100%;">
        <h2>
          VPN Profiles
          {{ if .User.Enabled }}
          <a href="/users/{{ .User.Id }}/profile_add" class="btn btn-sm btn-primary vertical-align-top pull-right">
            <i class="ion ion-android-add-circle"></i> Add Profile
          </a>
          {{ else }}
          <button type="button" class="btn btn-sm btn-default vertical-align-top pull-right" disabled>
            <i class="ion ion-android-add-circle"></i> Add Profile
          </button>
          {{ end }}
        </h2>
      </div>
    </div>

    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_content">
            <h4>
              Total Data Transfer: &nbsp; <i class="ion ion-ios-cloud-download-outline"></i>&nbsp;{{.TotalDownload}} &nbsp;<i class="ion ion-ios-cloud-upload-outline"></i>&nbsp; {{.TotalUpload}}
            </h4>
          </div>
        </div>
      </div>
    </div>

    <div class="clearfix"></div>

    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_content">
            <div class="table-responsive">
              <table class="table table-striped">
                <thead>
                  <tr>
                    <th>Description</th>
                    <th>User Name</th>
                    <th>Status</th>
                    <th>Login Count</th>
                    <th>Last Login</th>
                    <th>Inbound</th>
                    <th>Outbound</th>
                    <th> </th>
                  </tr>
                </thead>
                <tbody>
                  {{ if .ProfileList }}
                  {{ range .ProfileList }}
                  <tr id="tr_{{ .Id }}" class="clickable subspace-list-hover" data-toggle="collapse" data-target="#tr_target_{{ .Id }}">
                    <td>{{ .Description }}</td>
                    <td>{{ .Username }}</td>
                    <td>
                      {{ if eq .Enabled true }}
                      Enabled
                      {{ else }}
                      Disabled
                      {{ end }}
                    </td>
                    <td>{{ .LoginCount }}</td>
                    <td>{{ .LastLoginDate }}</td>
                    <td>{{ .IncomingBytes }}</td>
                    <td>{{ .OutgoingBytes }}</td>
                    <td><i class="glyphicon glyphicon-chevron-down"></i></td>
                  </tr>
                  <tr>
                    <td colspan="8" class="subspace-accordion-table-cell-style">
                      <div class="accordion-body collapse" id="tr_target_{{ .Id }}">
                        <div class="pull-right">
                          <button onclick="launchModal.call(this, 'Edit Profile', 'profileEditTemplate', 'profile_edit');" type="button" class="btn btn-primary btn-sm" data-profile-id="{{.Id}}"> <i class="ion ion-edit"></i> Edit</button>
                          {{ if .UserEnabled }}
                            {{ if eq .Enabled true }}
                            <button onclick="launchModal.call(this, 'Disable Profile', 'profileDisableTemplate', 'disable');" type="button" class="btn btn-warning btn-sm" data-profile-id="{{.Id}}"> <i class="ion ion-power"></i> Disable</button>
                            <button onclick="" type="button" class="btn btn-danger btn-sm" data-profile-id="{{.Id}}" disabled> <i class="ion ion-trash-a"></i> Delete</button>
                            {{ else }}
                            <button onclick="launchModal.call(this, 'Enable Profile', 'profileEnableTemplate', 'enable');" type="button" class="btn btn-success btn-sm" data-profile-id="{{.Id}}"> <i class="ion ion-power"></i> Enable</button>
                            <button onclick="launchModal.call(this, 'Delete Profile', 'profileDeleteTemplate', 'profile_delete');" type="button" class="btn btn-danger btn-sm" data-profile-id="{{.Id}}"> <i class="ion ion-trash-a"></i> Delete</button>
                            {{ end }}
                          {{ else }}
                          <button onclick="launchModal.call(this, 'Enable Profile', 'profileEnableTemplate', 'enable');" type="button" class="btn btn-success btn-sm" data-profile-id="{{.Id}}" disabled> <i class="ion ion-power"></i> Enable</button>
                          <button onclick="launchModal.call(this, 'Delete Profile', 'profileDeleteTemplate', 'profile_delete');" type="button" class="btn btn-danger btn-sm" data-profile-id="{{.Id}}"> <i class="ion ion-trash-a"></i> Delete</button>
                          {{ end }}
                        </div>
                        <table class="table table-condensed">
                          <thead>
                            <tr>
                              <th>Session Name</th>
                              <th>Source IP</th>
                              <th>Session Start Time</th>
                              <th>In</th>
                              <th>Out</th>
                              <th>Action</th>
                            </tr>
                          </thead>
                          <tbody>
                            {{ if .SessionList }}
                              {{ range .SessionList }}
                              <tr>
                                <td>{{ .SessionName }}</td>
                                <td>{{ .ClientIPAddress }}</td>
                                <td>{{ .ConnectionStartedAt }}</td>
                                <td>{{ .IncomingDataSize }}</td>
                                <td>{{ .OutgoingDataSize }}</td>
                                <td><button onclick="launchModal.call(this, 'Disconnect Session', 'profileDisconnecTemplate', 'session_delete');" type="button" class="btn btn-danger btn-sm" data-session-id="{{.SessionName}}"> <i class="glyphicon glyphicon-ban-circle"></i> </button></td>
                              </tr>
                              {{ end }}
                            {{ end }}
                          </tbody>
                        </table>
                      </div>
                    </td>
                  </tr>
                  {{ end }}
                {{ end }}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- /VPN Profile -->
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
          <h3 class="modal-title" id="modal-title">ACTION_TITLE</h3>
        </div>
        <div class="modal-body" id="modal-body">
          ACTION_CONTENT
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
          <!-- <a id="modal-action" href="ACTION" class="btn btn-danger" role="button"> Submit </a> -->
          <span id="modal-action"></span>
        </div>
      </form>
    </div>
  </div>
  <div id="user_id_refer" style="display: none;">{{ .User.Id }}</div>
</div>
<!-- /Dialogue -->
<!-- Dialogue content: edit profile -->
<div id="profileEditTemplate" style="display: none;">
  <label class="control-label col-md-3 col-sm-3 col-xs-12"> Description </label>
  <div class="col-md-9 col-sm-9 col-xs-12">
    <input id="edit-description" name="description" type="text" class="form-control" value="" placeholder="">
  </div>
  <div class="clearfix"></div>
  <br />
</div>
<!-- /Dialogue content: edit profile -->
<!-- Dialogue content: disable profile-->
<div id="profileDisableTemplate" style="display: none">
  <h4>Are you sure you want to do this?</h4>
  <div class="panel panel-warning">
    <div class="panel-heading">By performing the <b>Disable Profile</b> action:</div>
    <ul class="list-group">
      <li class="list-group-item">You will immediately disconnect all the VPN sessions associated with this Profile.</li>
      <li class="list-group-item">VPN connections using this Profile will be unusable.</li>
    </ul>
  </div>
</div>
<!-- /Dialogue content: disable profile -->
<!-- Dialogue content: enable profile-->
<div id="profileEnableTemplate" style="display: none">
  <h4>Are you sure you want to do this?</h4>
  <div class="panel panel-success">
    <div class="panel-heading">By performing the <b>Enable Profile</b> action:</div>
    <ul class="list-group">
      <li class="list-group-item">VPN connections using this Profile will be allowed.</li>
    </ul>
  </div>
</div>
<!-- /Dialogue content: enable user -->
<!-- Dialogue content: delete profile -->
<div id="profileDeleteTemplate" style="display: none;">
  <h4>Are you sure you want to do this?</h4>
  <div class="panel panel-danger">
    <div class="panel-heading">By performing the <b>Delete Profile</b> action:</div>
    <ul class="list-group">
      <li class="list-group-item">You will immediately disconnect all the VPN sessions associated with this Profile.</li>
      <li class="list-group-item">This Profile will be removed from this User and will be unusable.</li>
    </ul>
  </div>
</div>
<!-- /Dialogue content: delete profile -->
<!-- Dialogue content: disconnection -->
<div id="profileDisconnecTemplate" style="display: none;">
  <p>
    Are you sure you want to disconnect this session?
  </p>
  <br />
</div>
<!-- /Dialogue content: disconnection -->
