<div>
  <div class="row">
    <div class="col-md-8 col-sm-8 col-xs-12 col-md-offset-2 col-sm-offset-2">
      <div class="x_panel">
        <div class="x_title subspace-h1-title-visible">
          <h1>Setup Complete!</h1>
          <div class="clearfix"></div>
        </div>

        <div class="x_content">
          <form class="form-horizontal">
            <div class="form-group">
              <div class="col-md-1 col-md-offset-2 col-sm-1 col-xs-1">
                <span class="glyphicon glyphicon-ok-circle" style="font-size: 20px"></span>
              </div>
              <div class="col-md-8 col-sm-11 col-xs-11">
                <span style="font-size: 14px">Your Subspace Administrator account has been created.</span>
              </div>
            </div>

            <div class="form-group">
              <div class="col-md-1 col-md-offset-2 col-sm-1 col-xs-1">
                <span class="glyphicon glyphicon-ok-circle" style="font-size: 20px"></span>
              </div>
              <div id="downloadgroup" class="col-md-8 col-sm-11 col-xs-11 panel-group">
                <p style="font-size: 14px">
                  A VPN profile has been created for you.
                  <br /> Choose platform below for usage instructions.
                  <br />
                  <i>For security purposes, the downloadable profiles on this page will be unavailable after {{.vpnProfile.ttlInMinutes}} minutes.</i>
                </p>

                <!-- button accordion -->
                <div class="subspace-group-btn">
                  <button class="btn btn-primary" type="button" data-parent="#downloadgroup" data-toggle="collapse" data-target="#apple"><i class="ion ion-social-apple"></i> iOS / macOS</button>
                  <button class="btn btn-primary" type="button" data-parent="#downloadgroup" data-toggle="collapse" data-target="#windows"><i class="ion ion-social-windows"></i> Windows</button>
                  <button class="btn btn-primary" type="button" data-parent="#downloadgroup" data-toggle="collapse" data-target="#android"><i class="ion ion-social-android"></i> Android</button>
                </div>
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
                          <a class="btn btn-primary" href="{{.vpnProfile.appleProfilePath}}" download><i class="ion ion-archive"></i> Download</a>
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
                          <a class="btn btn-primary" href="{{.vpnProfile.windowsProfilePath}}" download><i class="ion ion-archive"></i> Download</a>
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

            <div class="form-group">
              <div class="col-md-1 col-md-offset-2 col-sm-1 col-xs-1">
                <span class="glyphicon glyphicon-ok-circle" style="font-size: 20px"></span>
              </div>
              <div class="col-md-8 col-sm-11 col-xs-11" style="font-size: 14px">
                When you're ready, continue on to the next page to enter the Subspace Dashboard for an overview of your service and management settings.
              </div>
            </div>

            <div class="ln_solid"></div>
            <div class="pull-right">
              <button type="button" class="btn btn-success" data-toggle="modal" data-target=".ss-modal">Complete Setup</button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</div>


<!-- Confirmation Dialog -->
<div class="modal fade ss-modal" tabindex="-1" role="dialog" aria-hidden="true">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-body">
        <div class="alert alert-danger" role="alert" style="font-size: 15px">
          Warning!
        </div>
        <p style="font-size: 14px; padding: 0px 8px;">
          You will not be able to view your VPN profile connection settings after leaving this page.
          <br /><br /> If necessary, press <b>Cancel</b> to go back.
        </p>
      </div>
      <div class="modal-footer">
        <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
        <a class="btn btn-success" style="color: #fff;" href="/" data-disable-with="Continue to Dashboard">Continue to Dashboard</a>
      </div>
    </div>
  </div>
</div>
