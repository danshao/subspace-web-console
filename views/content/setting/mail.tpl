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
        <h1>Mail Settings</h1>
      </div>
    </div>

    <div class="row">
      <div class="col-md-12 ">
        <div class="x_panel">
          <div class="x_title">
            <h2>SMTP Server</h2>
            <div class="clearfix"></div>
          </div>
          <div class="x_content">
            <div class="well">
              Subspace VPN profile connection settings can be sent directly to each user via email through your own SMTP server.
            </div>

            <form method="POST" class="form-horizontal">
              <div class="form-group">
                <label class="control-label col-md-3 col-sm-3 col-xs-12">SMTP Server</label>
                <div class="col-md-6 col-sm-6 col-xs-12">
                  <input id="smtpserver" name="smtpserver" type="text" class="form-control" value="{{.SystemInfo.SmtpHost}}">
                </div>
              </div>

              <div class="form-group">
                <label class="control-label col-md-3 col-sm-3 col-xs-12">SMTP Port</label>
                <div class="col-md-6 col-sm-6 col-xs-12">
                  <input id="smtpport" name="smtpport" type="text" class="form-control" value="{{ if eq .SystemInfo.SmtpPort 0 }}{{else}}{{.SystemInfo.SmtpPort}}{{ end }}">
                </div>
              </div>

              <div class="form-group">
                <label class="control-label col-md-3 col-sm-3 col-xs-12">Sender Name</label>
                <div class="col-md-6 col-sm-6 col-xs-12">
                  <input id="sendername" name="sendername" type="text" class="form-control" value="{{.SystemInfo.SmtpSenderName}}">
                </div>
              </div>

              <div class="form-group">
                <label class="control-label col-md-3 col-sm-3 col-xs-12">Sender Email</label>
                <div class="col-md-6 col-sm-6 col-xs-12">
                  <input id="senderemail" name="senderemail" type="text" class="form-control" value="{{.SystemInfo.SmtpSenderEmail}}">
                </div>
              </div>

              <div class="form-group">
                <div class="col-md-6 col-sm-6 col-xs-12 col-md-offset-3 col-sm-offset-3">
                  <div class="checkbox">
                    <label class="label-cr">
                      {{ if eq .SystemInfo.SmtpAuthentication false }}
                      <input id="authentication" name="authentication" type="checkbox">
                      {{ else }}
                      <input id="authentication" name="authentication" type="checkbox" checked/>
                      {{ end }}
                      <span class="labelbox"><i class="ico"></i>Authentication required</span>
                    </label>
                  </div>
                </div>
              </div>

              <div class="form-group">
                <label class="control-label col-md-3 col-sm-3 col-xs-12">Username</label>
                <div class="col-md-6 col-sm-6 col-xs-12">
                  {{ if eq .SystemInfo.SmtpAuthentication false }}
                  <input id="username" name="username" type="text" class="form-control" value="" disabled>
                  {{ else }}
                  <input id="username" name="username" type="text" class="form-control" value="{{.SystemInfo.SmtpUsername}}">
                  {{ end }}
                </div>
              </div>

              <div class="form-group">
                <label class="control-label col-md-3 col-sm-3 col-xs-12">Password</label>
                <div class="col-md-6 col-sm-6 col-xs-12">
                  {{ if eq .SystemInfo.SmtpAuthentication false }}
                  <input id="password" name="password" type="password" class="form-control" value="" disabled>
                  {{ else }}
                  <input id="password" name="password" type="password" class="form-control" value="{{.SystemInfo.SmtpPassword}}">
                  {{ end }}
                </div>
              </div>

              <div class="form-group pull-right">
                <div class="col-md-12 text-right">
                  <button type="submit" class="btn btn-primary" data-disable-with="Save Settings">Save Settings</button>
                </div>
              </div>

            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
