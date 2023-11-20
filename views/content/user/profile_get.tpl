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

    <div class="page-title">
      <div class="title_left">
        <h1>VPN Profile Information</h1>
      </div>
    </div>

    <div class="clearfix"></div>
    <div class="ln_solid"></div>
    <div class="row">
      <div class="col-md-12">
        <div class="form-group">
          <div id="downloadgroup" class="col-md-10 panel-group">
            <p style="font-size: 14px">
              A new VPN Profile <b>{{.vpnProfile.description}}</b> has been created for <b>{{.User.Email}}</b>.
              <br /><br />
              Choose platform below for usage instructions.
              <br />
              <i>For security purposes, the downloadable profiles on this page will be unavailable after {{.vpnProfile.ttlInMinutes}} minutes.</i>
            </p>

            <!-- button accordion -->
            <button class="btn btn-primary" type="button" data-parent="#downloadgroup" data-toggle="collapse" data-target="#apple"><i class="ion ion-social-apple"></i> iOS / macOS</button>
            <button class="btn btn-primary" type="button" data-parent="#downloadgroup" data-toggle="collapse" data-target="#windows"><i class="ion ion-social-windows"></i> Windows</button>
            <button class="btn btn-primary" type="button" data-parent="#downloadgroup" data-toggle="collapse" data-target="#android"><i class="ion ion-social-android"></i> Android</button>
            <!-- /button accordion-->

            <div class="clearfix"></div>
            <br />

            <!-- APPLE -->
            <div class="panel">
              <div class="collapse" id="apple">
                <div class="well">
                  <i class="ion ion-social-apple" style="font-size: 18px"></i>
                  <p style="padding-top: 10px">
                    Please click the <b>Download</b> button. Locate the downloaded profile and open it. The profile will automatically install onto your device.
                  </p>
                  <div class="clearfix">
                    <div class="pull-right">
                      <a class="btn btn-primary " href="{{.vpnProfile.appleProfilePath}}" download><i class="ion ion-archive"></i> Download</a>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- WINDOWS -->
            <div class="panel">
              <div class="collapse" id="windows">
                <div class="well">
                  <i class="ion ion-social-windows" style="font-size: 18px"></i>
                  <p style="padding-top: 10px">
                    Please click the <b>Download</b> button. Locate the downloaded profile and open it. The profile will automatically install onto your device.
                  </p>
                  <p>
                    Additionally, you will need the information below to complete installation. Be sure to write it down now!
                  </p>
                  <div class="panel panel-primary">
                          <div class="panel-heading">Profile Information</div>
                          <table class="table">
                            <tr>
                              <td style="width: 30%" align="center">Username</td>
                              <td>{{.vpnProfile.username}}</td>
                            </tr>
                            <tr>
                              <td style="width: 30%" align="center">Password</td>
                              <td>{{.vpnProfile.password}}</td>
                            </tr>
                          </table>
                        </div>
                        <div class="panel panel-primary">
                          <div class="panel-heading">Connection Information</div>
                          <table class="table">
                            <tr>
                              <td style="width: 30%" align="center">Host</td>
                              <td>{{.vpnProfile.host}}</td>
                            </tr>
                            <tr>
                              <td style="width: 30%" align="center">Pre-Shared Key</td>
                              <td>{{.vpnProfile.key}}</td>
                            </tr>
                          </table>
                        </div>
                  <div class="clearfix">
                    <div class="pull-right">
                      <a class="btn btn-primary " href="{{.vpnProfile.windowsProfilePath}}" download><i class="ion ion-archive"></i> Download</a>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- ANDROID -->
            <div class="panel">
              <div class="collapse" id="android">
                <div class="well">
                  <i class="ion ion-social-android" style="font-size: 18px"></i>
                  <p style="padding-top: 10px">
                    Record the information below before leaving this page. You will need to set up the VPN Profile manually on your device.
                  </p>
                  <div class="panel panel-primary">
                    <div class="panel-heading">Profile Information</div>
                    <table class="table">
                      <tr>
                        <td style="width: 30%" align="center">Username</td>
                        <td>{{.vpnProfile.username}}</td>
                      </tr>
                      <tr>
                        <td style="width: 30%" align="center">Password</td>
                        <td>{{.vpnProfile.password}}</td>
                      </tr>
                    </table>
                  </div>
                  <div class="panel panel-primary">
                    <div class="panel-heading">Connection Information</div>
                    <table class="table">
                      <tr>
                        <td style="width: 30%" align="center">Host</td>
                        <td>{{.vpnProfile.host}}</td>
                      </tr>
                      <tr>
                        <td style="width: 30%" align="center">Pre-Shared Key</td>
                        <td>{{.vpnProfile.key}}</td>
                      </tr>
                    </table>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="text-right">
      <a href="/users/{{.User.Id}}" class="btn btn-primary">Back to User Info</a>
    </div>
  </div>
</div>
