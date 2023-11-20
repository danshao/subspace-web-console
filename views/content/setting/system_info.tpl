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

    <!-- Main Content -->
    <div class="page-title">
      <div class="title_left">
        <h1>System Information</h1>
      </div>
    </div>
    <div class="bs-callout bs-callout-danger">
      <h4>Attention</h4>
      <p>Modifications to these settings can possibly cause your users to be unable to connect to the Subspace service if corresponding actions are not taken by your users. It is your responsibility to maintain a consistent experience by helping your users properly react to the changes here.</p>
    </div>
    <!-- Host info -->
    <div class="row">
      <div class="col-md-12 col-sm-12 col-xs-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>Host</h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <!-- form -->
            <form>
              <div class="form-group">
                <label>DNS Hostname &nbsp; <small class="subspace-label-description"> updated at: {{ localTimeFmt .SystemInfo.HostUpdatedDate }} </small></label>
                <div class="clearfix">
                  <p id="hostname" class="form-control-static" style="float: left;">{{ if eq .SystemInfo.Host "" }} Not set {{ else }} {{ .SystemInfo.Host }} {{ end }}</p>
                  <div class="btn-group" style="float: right;">
                    <button onclick="verifyDNS.call(this, {{.SystemInfo.Host}});" id="dns-verify" class="btn btn-default" type="button" {{ if eq "" .SystemInfo.Host }} disabled {{ end }}>
                      <i class="fa fa-check-circle-o" style="display:none;"></i>
                      <i class="fa fa-times-circle-o" style="display:none;"></i>
                      <i class="fa fa-spinner fa-pulse fa-fw" style="display:none;"></i>
                      Verify
                    </button>
                    <button onclick="launchModal.call(this, 'Update DNS Hostname', 'hostnameEditTemplate', 'hostname_edit');" type="button" class="btn btn-primary" data-id="">
                      <i class="ion ion-edit"></i> Edit
                    </button>
                  </div>
                  <span class="refer" style="display:none;">{{ .SystemInfo.Host }}</span>
                </div>
                <font id="host-verify-message" color="red" style="display:none;"></font>
              </div>
              <br>
              <div class="form-group">
                <label>Current IP Address</label>
                <p id="ipaddress" class="form-control-static">{{ if eq .SystemInfo.Ip "" }} Not set {{ else }} {{ .SystemInfo.Ip }} {{ end }}</p>
              </div>
            </form>
            <!-- /form -->
          </div>
        </div>
      </div>
    </div>
    <!-- /Host info -->

    <div class="ln_solid"></div>
    <div class="clearfix"></div>

    <!-- VPN info -->
    <div class="row">
      <div class="col-md-12">
        <div class="x_panel">
          <div class="x_title">
            <h2>VPN</h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <!-- form -->
            <form>
              <div class="form-group">
                <label>Pre-shared Key &nbsp; <small class="subspace-label-description"> updated at: {{ localTimeFmt .SystemInfo.PreSharedKeyUpdatedDate }} </small></label>
                <div class="clearfix">
                  <p id="presharedkey" class="form-control-static" style="float: left;">{{ if eq .SystemInfo.PreSharedKey "" }} Not set {{ else }} {{ .SystemInfo.PreSharedKey }} {{ end }}</p>
                  <button onclick="launchModal.call(this, 'Update Pre-shared Key', 'presharedKeyEditTemplate', 'presharedkey_edit');" type="button" class="btn btn-primary" style="float: right;" data-id=""> <i class="ion ion-edit"></i> Edit</button>
                  <span class="refer" style="display:none;">{{ .SystemInfo.PreSharedKey }}</span>
                </div>
              </div>
              <br>
              <div class="form-group">
                <label>UUID &nbsp; <small class="subspace-label-description"> updated at: {{ localTimeFmt .SystemInfo.UuidUpdatedDate }} </small></label>
                <div class="clearfix">
                  <p id="uuid" class="form-control-static" style="float: left;">{{ if eq .SystemInfo.Uuid "" }} Not set {{ else }} {{ .SystemInfo.Uuid }} {{ end }}</p>
                  <button onclick="launchModal.call(this, 'Update UUID', 'uuidEditTemplate', 'uuid_edit');" type="button" class="btn btn-primary" style="float: right;" data-id=""> <i class="ion ion-edit"></i> Edit</button>
                  <span class="refer" style="display:none;">{{ .SystemInfo.Uuid }}</span>
                </div>
              </div>
            </form>
            <!-- /form -->
          </div>
        </div>
      </div>
    </div>
    <!-- /VPN info -->
    <!-- /Main Content -->
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
          <button type="button" class="close" data-dismiss="modal"><span aria-hidden="true">×</span></button>
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

<!-- Dialogue content: edit DNS host name -->
<div id="hostnameEditTemplate" style="display: none;">
  <div class="panel panel-danger">
    <div class="panel-heading">
      <div class="panel-title">
        Warning
      </div>
    </div>
    <ul class="list-group">
      <li class="list-group-item">
        This setting is system-level configuration, changing this value may affect all VPN connections. Please make sure that your users are informed about these changes before saving them.
      </li>
      <li class="list-group-item">
        If you have properly configured your SMTP settings, an email with the updated changes will be sent out to each user.
      </li>
    </ul>
  </div>
  <div class="clearfix"></div>
  <br />
  <label class="control-label col-md-3 col-sm-3 col-xs-12"> DNS Hostname </label>
  <div class="col-md-9 col-sm-9 col-xs-12">
    <input name="hostname" type="text" class="form-control edit-description" value="" placeholder="">
    <font class="error-msg" color="red;"></font>
  </div>
  <div class="clearfix"></div>
  <br />
</div>
<!-- /Dialogue content: edit DNS host name -->

<!-- Dialogue content: edit pre-shared key -->
<div id="presharedKeyEditTemplate" style="display: none;">
  <div class="panel panel-danger">
    <div class="panel-heading">
      <div class="panel-title">
        Warning
      </div>
    </div>
    <ul class="list-group">
      <li class="list-group-item">
        This setting is system-level configuration, changing this value may affect all VPN connections. Please make sure that your users are informed about these changes before saving them.
      </li>
      <li class="list-group-item">
        If you have properly configured your SMTP settings, an email with the updated changes will be sent out to each user.
      </li>
    </ul>
  </div>
  <div class="clearfix"></div>
  <br />
  <label class="control-label col-md-3 col-sm-3 col-xs-12"> Pre-shared Key </label>
  <div class="col-md-9 col-sm-9 col-xs-12">
    <input name="presharedkey" type="text" class="form-control edit-description" value="" placeholder="">
    <font class="error-msg" color="red;"></font>
  </div>
  <div class="clearfix"></div>
  <br />
</div>
<!-- /Dialogue content: edit pre-shared key -->

<!-- Dialogue content: edit UUID -->
<div id="uuidEditTemplate" style="display: none;">
  <!--<label class="control-label col-md-3 col-sm-3 col-xs-12"> UUID </label>-->
  <div class="col-md-12">
    <p>UUID will be updated after pressing <strong>re-generate</strong> button.</p>
    <!--<input name="uuid" type="text" class="form-control edit-description" value="" placeholder="">
    <font class="error-msg" color="red;"></font>-->
  </div>
  <div class="clearfix"></div>
  <br />
</div>
<!-- /Dialogue content: edit UUID -->
